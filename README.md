# file-reader

## Running locally (for AVRO)

1. Install homebrew `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
2. Install awscli `brew install awscli`
3. Installed gogen-avro to generate avro `go install github.com/actgardner/gogen-avro/v10/cmd/...@latest`
4. Create your file
5. Installed kcat to read avro messages from topic `brew install kcat`

Start up redpanda and localstack via docker-compose
`docker-compose up -d`

Ensure the env variables have been exported via `direnv`.

Build avros
`sh ./src/scripts/build-avros.sh`

Create topics
`sh ./src/scripts/create-topics.sh`

Run service
`go run ./src/services/organization/server`

Put S3 file on localstack and put file create message on kafka topic to trigger the ingest
(ToDo: write integration test for this)

Can get avro messages from topic manually
`kcat -b localhost:9092 -t organization  -s avro  -r http://localhost:8081`

## Lambda

You'll require the S3FileCreatedUpdated topic to be created
```
docker exec redpanda-1 rpk topic create S3FileCreatedUpdated --brokers=localhost:9092
```

### deploy lambda and listen to org bucket

```
docker-compose up -d && \
sleep 1 && \
docker exec redpanda-1 rpk topic create S3FileCreatedUpdated --brokers=localhost:9092 && \
aws --endpoint-url http://localhost:4566 lambda create-function --function-name fileevent --handler main --runtime go1.x --role your-role --zip-file fileb://main.zip --environment 'Variables={KAFKA_BROKER=redpanda:29092,SCHEMA_REGISTRY=http://redpanda:8081,TOPIC_NAME=S3FileCreatedUpdated}' && \
aws --endpoint-url=http://localhost:4566 s3api put-bucket-notification-configuration --bucket organization --notification-configuration file://scripts/local/s3-notif-config.json
```

### teardown of lambda
```
aws --endpoint-url http://localhost:4566 lambda delete-function --function-name fileevent
``` 