package kafka

import (
	"context"
	"file_reader/src/config"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"go.uber.org/zap"
)

type OrganizationProducer interface {
	PublishCreate(ctx context.Context, msgs ...kafka.Message) error
	PublishUpdate(ctx context.Context, msgs ...kafka.Message) error
	Close()
	Run()
	GetNewKafkaWriter(topic string) *kafka.Writer
}

type organizationProducer struct {
	log          zap.Logger
	cfg          *config.Config
	createWriter *kafka.Writer
	updateWriter *kafka.Writer
}

// contructor

func NewOrganizationProducer(log zap.Logger, cfg *config.Config) *organizationProducer {
	return &organizationProducer{log: log, cfg: cfg}
}

// Create new kafka writer

func (o *organizationProducer) GetNewKafkaWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(o.cfg.Kafka.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(o.log.Sugar().Debugf),
		ErrorLogger:  kafka.LoggerFunc(o.log.Sugar().Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
	return w
}

// Run init producer writers
func (o *organizationProducer) Run() {
	o.createWriter = o.GetNewKafkaWriter(createOrganizationTopic)
	o.updateWriter = o.GetNewKafkaWriter(updateOrganizationTopic)
}

// Close close writers
func (o organizationProducer) Close() {
	o.createWriter.Close()
	o.updateWriter.Close()
}

// Publish message to create topic
func (o *organizationProducer) PublishCreate(ctx context.Context, msgs ...kafka.Message) error {
	return o.createWriter.WriteMessages(ctx, msgs...)
}

// Publish message to update topic

func (o *organizationProducer) PublishUpdate(ctx context.Context, msgs ...kafka.Message) error {
	return o.updateWriter.WriteMessages(ctx, msgs...)
}
