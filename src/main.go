package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "dummy_topic"
	brokerAddress = "localhost:9092"
	csvPath       = "data.csv"
)

func createKafkaWriter(topic string, brokerAddr string, logger *log.Logger) *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		// assign the logger to the writer
		Logger: logger,
	})
	return w
}

func ingestCsvToKafka(f *os.File, w *kafka.Writer, ctx context.Context) {
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", rec)
		msg := strings.Join(rec, "-")
		// Write the message to the Kafka topic
		// This is slow because it isn't batching messages together
		err = w.WriteMessages(ctx, kafka.Message{
			Key: []byte(msg),
			// create an arbitrary message payload for the value
			Value: []byte("this is message:" + msg),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}
	}
}

func main() {
	l := log.New(os.Stdout, "kafka writer: ", 0)
	ctx := context.Background()

	w := createKafkaWriter(topic, brokerAddress, l)

	// open file
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	ingestCsvToKafka(f, w, ctx)
}
