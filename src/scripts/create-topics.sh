#!/bin/sh

docker exec -it redpanda-1 rpk topic create person --brokers=localhost:9092