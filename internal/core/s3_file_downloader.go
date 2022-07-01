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
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Reader interface {
	Read() ([]string, error)
}

type Client struct {
	s3iface.S3API
	key    string
	bucket string
}
type Filter struct {
	k, v string
}

func NewClient(s s3iface.S3API, bucket, key string) *Client {
	return &Client{
		S3API:  s,
		bucket: bucket,
		key:    key,
	}
}

// NewFilter creates a new filter foa key value pair
func NewFilter(k, v string) Filter {
	return Filter{k, v}
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

	// Process the headers first
	// Read the rest of the first chunk
	var lines [][]byte
	var colNum int

	// If the chunk does not contain the whole header, keep trying with larger chunk size
	for {
		l, c, newChunkSize, err := processCsvHeaders(ctx, logger, downloader, s3FileCreated, chunkSize, fileRows)
		if err != nil {
			return err
		}
		lines = l
		colNum = c
		if newChunkSize == chunkSize {
			break
		}
		chunkSize = newChunkSize

	}

	// Read the remaining chunks
	partialLine := bytes.Trim(lines[len(lines)-1], "\x00")
	start := chunkSize + 1
	for start < int(s3FileCreated.Payload.Content_length) {
		buff, err := downloadChunk(ctx, logger, downloader, s3FileCreated, start, chunkSize)
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

func processCsvHeaders(ctx context.Context, logger *zaplogger.ZapLogger, downloader *s3manager.Downloader, s3FileCreated avro.S3FileCreatedUpdated, chunkSize int, fileRows chan []string) ([][]byte, int, int, error) {
	buff, err := downloadChunk(ctx, logger, downloader, s3FileCreated, 0, chunkSize)
	if err != nil {
		return nil, 0, chunkSize, err
	}

	lines := bytes.Split(buff, []byte("\n"))
	// If the chunk does not contain the whole headers, double the chunk size
	if len(lines) == 1 {
		return nil, 0, chunkSize * 2, nil
	}

	headers := strings.Split(string(lines[0]), ",")
	colNum := len(headers)
	fileRows <- headers

	for line := 1; line < len(lines)-1; line++ {
		cols := strings.Split(string(lines[line]), ",")
		if len(cols) != colNum {
			logger.Error(ctx, "Error streaming from file. Wrong number of cols", cols)
		}
		fileRows <- cols
	}
	return lines, colNum, chunkSize, nil
}

func downloadChunk(
	ctx context.Context,
	logger *zaplogger.ZapLogger,
	downloader *s3manager.Downloader,
	s3FileCreated avro.S3FileCreatedUpdated,
	start int,
	chunkSize int,
) ([]byte, error) {
	buff := make([]byte, chunkSize)
	awsbuff := aws.NewWriteAtBuffer(buff)
	numBytes, err := downloader.DownloadWithContext(ctx, awsbuff,
		&s3.GetObjectInput{
			Bucket: aws.String(s3FileCreated.Payload.Bucket_name),
			Key:    aws.String(s3FileCreated.Payload.Key),
			Range:  aws.String(fmt.Sprintf("bytes=%v-%v", start, start+chunkSize)),
		})
	if err != nil {
		return nil, err
	}
	logger.Infof(ctx, "Streamed %v bytes from %s. Range(%v-%v)", numBytes, s3FileCreated.Payload.Key, start, start+chunkSize)
	return buff, nil
}
