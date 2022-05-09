#!/bin/sh

docker exec -it redpanda-1 rpk topic create organization --brokers=localhost:9092
docker exec -it redpanda-1 rpk topic create s3filecreation --brokers=localhost:9092