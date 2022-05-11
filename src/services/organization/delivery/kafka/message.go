package kafka

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// TopicPartition is a generic placeholder for a Topic+Partition and optionally Offset.
type TopicPartition struct {
	Topic     *string
	Partition int32
	Offset    int
	Metadata  *string
	Error     error
}

// Message represents a Kafka message
type Message struct {
	TopicPartition TopicPartition
	Value          []byte
	Key            []byte
	Timestamp      time.Time
	Opaque         interface{}
	Headers        []kafka.Header
}

// String returns a human readable representation of a Message.
// Key and payload are not represented.
func (m *Message) String() string {
	var topic string
	if m.TopicPartition.Topic != nil {
		topic = *m.TopicPartition.Topic
	} else {
		topic = ""
	}
	return fmt.Sprintf("%s[%d]@%s", topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
}
