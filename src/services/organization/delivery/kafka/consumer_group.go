package kafka

import (
	"context"
	"encoding/json"
	"file_reader/src/config"
	"file_reader/src/models"
	"file_reader/src/services/organization"
	"sync"

	"github.com/go-playground/validator"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"go.uber.org/zap"
)

type OrganizationsConsumerGroup struct {
	Brokers        []string
	GroupID        string
	log            zap.Logger
	cfg            *config.Config
	organizationUC organization.UseCase
	validate       *validator.Validate
}

// NewOrganizationsConsumerGroup constructor
func NewOrganizationsConsumerGroup(
	brokers []string,
	groupId string,
	log zap.Logger,
	cfg *config.Config,
	organizationUC organization.UseCase,
	validate *validator.Validate,

) *OrganizationsConsumerGroup {
	return &OrganizationsConsumerGroup{Brokers: brokers, GroupID: groupId, log: log, cfg: cfg, organizationUC: organizationUC, validate: validate}
}

func (ocg *OrganizationsConsumerGroup) getNewKafkaReader(kafkaURL []string, topic, groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupId,
		Topic:                  topic,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: partitionWatchInterval,
		Logger:                 kafka.LoggerFunc(ocg.log.Sugar().Debugf),
		ErrorLogger:            kafka.LoggerFunc(ocg.log.Sugar().Errorf),
		MaxAttempts:            maxAttempts,
		Dialer: &kafka.Dialer{
			Timeout: dialTimeout,
		},
	})
}

func (ocg *OrganizationsConsumerGroup) getNewKafkaWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(ocg.Brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Logger:       kafka.LoggerFunc(ocg.log.Sugar().Debugf),
		ErrorLogger:  kafka.LoggerFunc(ocg.log.Sugar().Errorf),
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}

	return w

}

func (ocg *OrganizationsConsumerGroup) consumeCreateOrganization(
	ctx context.Context,
	cancel context.CancelFunc,
	groupId string,
	topic string,
	workersNum int,
) {
	r := ocg.getNewKafkaReader(ocg.Brokers, topic, groupId)

	defer cancel()

	defer func() {
		if err := r.Close(); err != nil {
			ocg.log.Sugar().Errorf("r.Close", zap.String("error", err.Error()))
			cancel()
		}
	}()

	w := ocg.getNewKafkaWriter(deadLetterQueueTopic)

	defer func() {
		if err := w.Close(); err != nil {
			ocg.log.Error("w.Close", zap.String("error", err.Error()))
			cancel()
		}
	}()

	ocg.log.Sugar().Infof("Starting consumer group: %v", r.Config().GroupID)

	wg := &sync.WaitGroup{}

	for i := 0; i <= workersNum; i++ {
		wg.Add(1)
		go ocg.createOrganizationWorker(ctx, cancel, r, w, wg, i)
	}
	wg.Wait()

}

func (ocg *OrganizationsConsumerGroup) publishErrorMessage(ctx context.Context, w *kafka.Writer, m kafka.Message, err error) error {
	errMsg := &models.ErrorMessage{
		Offset:    m.Offset,
		Error:     err.Error(),
		Time:      m.Time.UTC(),
		Partition: m.Partition,
		Topic:     m.Topic,
	}

	errMsgBytes, err := json.Marshal(errMsg)

	if err != nil {
		return err
	}

	return w.WriteMessages(ctx, kafka.Message{
		Value: errMsgBytes,
	})
}

func (ocg *OrganizationsConsumerGroup) RunConsumers(ctx context.Context, cancel context.CancelFunc) {
	go ocg.consumeCreateOrganization(ctx, cancel, organizationsGroupID, createOrganizationTopic, createOrganizationWorkers)
}
