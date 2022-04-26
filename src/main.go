package main

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:9092"
	csvPath       = "data.csv"
)

// topic, operation and rowToOperation could be defined for each operation type, then we compose the reader out of them
const topic = "organisation"
const schemaPath = "./src/avros/organisation.avsc"

type metadata struct {
	Origin_application string
	Region             string
	Tracking_id        string
}
type payload struct {
	Organization_name string
	Guid              string
}
type organisation struct {
	Payload  payload
	Metadata metadata
}

func rowToOrganisation(row []string) organisation {
	md := metadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        uuid.NewString(),
	}
	pl := payload{
		Guid:              row[0],
		Organization_name: row[1],
	}
	return organisation{Payload: pl, Metadata: md}
}

func composeCsvIngester[operation any](
	topic string,
	w *kafka.Writer,
	rowToOpConverter func([]string) operation,
	schema *srclient.Schema,
	brokerAddrs []string,
	logger *log.Logger,
) func(f *os.File, ctx context.Context) {
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	return func(f *os.File, ctx context.Context) {
		csvReader := csv.NewReader(f)
		for {
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			// Map row to bytes using schema
			op := rowToOpConverter(row)
			value, _ := json.Marshal(op)
			native, _, _ := schema.Codec().NativeFromTextual(value)
			valueBytes, err := schema.Codec().BinaryFromNative(nil, native)
			if err != nil {
				panic("error converting to avro bytes " + err.Error())
			}

			// Combine row bytes with schema id to make a record
			var recordValue []byte
			recordValue = append(recordValue, byte(0))
			recordValue = append(recordValue, schemaIDBytes...)
			recordValue = append(recordValue, valueBytes...)

			// Put the row on the topic
			key, _ := uuid.NewUUID()
			err = w.WriteMessages(
				ctx,
				kafka.Message{
					Key:   []byte(key.String()),
					Value: recordValue,
				},
			)
			if err != nil {
				panic("could not write message " + err.Error())
			}
		}
	}
}

// Just using this to test that the serialize/deserialize is working
func printTopicMessages(topic string, ctx context.Context, n int) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	schemaRegistryClient := srclient.CreateSchemaRegistryClient("http://localhost:8081")

	for i := 1; i < n; i++ {
		msg, err := r.ReadMessage(ctx)
		if err == nil {
			schemaID := binary.BigEndian.Uint32(msg.Value[1:5])
			schema, err := schemaRegistryClient.GetSchema(int(schemaID))
			if err != nil {
				panic(fmt.Sprintf("Error getting the schema with id '%d' %s", schemaID, err))
			}
			native, _, _ := schema.Codec().NativeFromBinary(msg.Value[5:])
			value, _ := schema.Codec().TextualFromNative(nil, native)
			fmt.Printf("Here is the message %s\n", string(value))
		} else {
			fmt.Printf("Error consuming the message: %v (%v)\n", err, msg)
			break
		}
	}
}

func main() {
	//Setup
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Get schema
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}
<<<<<<< Updated upstream
=======
	schemaRegistryClient := srclient.CreateSchemaRegistryClient("http://localhost:8081")
	schema, err := schemaRegistryClient.CreateSchema(topic, string(schemaBytes), "AVRO")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Schema '%d' retrieved successfully!\n", schema.ID())
>>>>>>> Stashed changes

	logger := log.New(os.Stdout, "kafka writer: ", 0)
	ctx := context.Background()

<<<<<<< Updated upstream
	// Compose the csv ingester
=======
	brokerAddrs := []string{brokerAddress}
	fmt.Printf("Broker: %s\n", brokerAddrs[0])
	// Prepare kafka producer for the provided topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokerAddrs,
		Topic:   topic,
		Logger:  logger,
	})

>>>>>>> Stashed changes
	csvIngester := composeCsvIngester(
		topic,
		w,
		rowToOrganisation,
		schema,
		[]string{brokerAddress},
		logger,
	)

	// Actually put the csv rows onto the topic
	csvIngester(f, ctx)

	// Consume from the topic to prove it worked
	printTopicMessages(topic, ctx, 5)

}
