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

