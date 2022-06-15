#!/bin/sh

docker exec -it redpanda-1 rpk topic create organization-membership-avro --brokers=localhost:9092
docker exec -it redpanda-1 rpk topic create S3FileCreatedUpdated --brokers=localhost:9092
docker exec -it redpanda-1 rpk topic create organization-proto --brokers=localhost:9092
