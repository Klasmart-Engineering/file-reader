package core

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	zaplogger "github.com/KL-Engineering/file-reader/internal/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Reader interface {
	Read() ([]string, error)
}

func DownloadFile(ctx context.Context, logger *zaplogger.ZapLogger, awsSession *session.Session, s3FileCreated avro.S3FileCreated, fileRows chan []string) error {
	// Download file from S3 and pass rows into a channel of []string

	// Create and open file on /tmp/
	f, err := ioutil.TempFile("", "file-reader-"+s3FileCreated.Payload.Key)
	if err != nil {
		log.Fatal("Failed to make tmp file", err)
	}
	defer os.Remove(f.Name())

	// Download from S3 to file
	downloader := s3manager.NewDownloader(awsSession)
	numBytes, err := downloader.Download(f,
		&s3.GetObjectInput{
			Bucket: aws.String(s3FileCreated.Payload.Bucket_name),
			Key:    aws.String(s3FileCreated.Payload.Key),
		})
	if err != nil {
		return err
	}
	logger.Infof(ctx, "Downloaded %s %d bytes", s3FileCreated.Payload.Key, numBytes)

	// Close and reopen the same file for ingest (until thought of alternative)
	f.Close()
	f, _ = os.Open(f.Name())

	// Read file and pass rows to fileRows channel
	go ReadRows(ctx, logger, f, s3FileCreated.Payload.Content_type, fileRows)

	return nil
}
