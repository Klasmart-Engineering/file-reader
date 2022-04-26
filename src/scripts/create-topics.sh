#!/bin/sh

docker exec -it redpanda-1 rpk topic create organisation --brokers=localhost:9092