package s3event

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type S3Event struct {
	message kafka.Message
}

// UseCase Product
type UseCase interface {
	processWithDownload(ctx context.Context, s3Event S3Event) error
	processWithStreaming(ctx context.Context, s3Event S3Event) error
}
