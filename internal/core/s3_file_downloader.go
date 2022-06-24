package core

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func DownloadFile(ctx context.Context, logger *zaplogger.ZapLogger, awsSession *session.Session, s3FileCreated avro.S3FileCreatedUpdated, fileRows chan []string) error {
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

func StreamFile(ctx context.Context, logger *zaplogger.ZapLogger, awsSession *session.Session, s3FileCreated avro.S3FileCreatedUpdated, fileRows chan []string, chunkSize int) error {
	defer func() {
		close(fileRows)
	}()

	if chunkSize < int(s3FileCreated.Payload.Content_length) {
		chunkSize = int(s3FileCreated.Payload.Content_length)
	}

	// Download the first chunk which includes the headers
	downloader := s3manager.NewDownloader(awsSession)
	buff := make([]byte, chunkSize)
	err := downloadChunk(ctx, logger, downloader, buff, s3FileCreated, 0, chunkSize)
	if err != nil {
		return err
	}

	// Read the headers first
	lines := bytes.Split(buff, []byte("\n"))
	headers := strings.Split(string(lines[0]), ",")
	colNum := len(headers)
	if err != nil {
		logger.Error(ctx, "Error reading headers when streaming file", err)
		return err
	}
	fileRows <- headers
	// Read the rest of the first chunk
	for line := 1; line < len(lines)-1; line++ {
		cols := strings.Split(string(lines[line]), ",")
		if len(cols) != colNum {
			logger.Error(ctx, "Error streaming from file. Wrong number of cols", cols)
		}
		fileRows <- cols
	}

	// Read the remaining chunks
	partialLine := bytes.Trim(lines[len(lines)-1], "\x00")
	start := chunkSize + 1
	for start < int(s3FileCreated.Payload.Content_length) {
		buff := make([]byte, chunkSize)
		err := downloadChunk(ctx, logger, downloader, buff, s3FileCreated, start, chunkSize)
		if err != nil {
			return err
		}

		lines := bytes.Split(buff, []byte("\n"))
		// Edge case for when the batch only contains a single partial line
		if len(lines) < 2 {
			partialLine = append(partialLine, bytes.Trim(buff, "\x00")...)
			start = start + chunkSize + 1
			continue
		}
		// Feed each line in this batch to the fileRows channel
		for line := 0; line < len(lines)-1; line++ {
			var l []byte
			if line == 0 {
				// Stitch first line in batch to the last line in the prev batch as they are the same line
				l = append(partialLine, lines[0]...)
			} else {
				l = lines[line]
			}
			logger.Info(ctx, "line size in bytes", len(l))
			cols := strings.Split(string(l), ",")
			if len(cols) != colNum {
				logger.Error(ctx, "Error streaming from file. Wrong number of cols", cols)
			}
			fileRows <- cols
		}
		partialLine = bytes.Trim(lines[len(lines)-1], "\x00")
		start = start + chunkSize + 1
	}
	if len(partialLine) > 0 {
		cols := strings.Split(string(partialLine), ",")
		if len(cols) != colNum {
			logger.Error(ctx, "Error streaming from file. Wrong number of cols", cols)
		}
		fileRows <- cols
	}

	return nil
}

func downloadChunk(
	ctx context.Context,
	logger *zaplogger.ZapLogger,
	downloader *s3manager.Downloader,
	buff []byte,
	s3FileCreated avro.S3FileCreatedUpdated,
	start int,
	chunkSize int,
) error {
	awsbuff := aws.NewWriteAtBuffer(buff)
	numBytes, err := downloader.DownloadWithContext(ctx, awsbuff,
		&s3.GetObjectInput{
			Bucket: aws.String(s3FileCreated.Payload.Bucket_name),
			Key:    aws.String(s3FileCreated.Payload.Key),
			Range:  aws.String(fmt.Sprintf("bytes=%v-%v", start, start+chunkSize)),
		})
	if err != nil {
		return err
	}
	logger.Infof(ctx, "Streamed %v bytes from %s. Range(%v-%v)", numBytes, s3FileCreated.Payload.Key, start, start+chunkSize)
	return nil
}
