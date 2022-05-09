#!/bin/sh

docker exec -it redpanda-1 rpk topic create organization_proto --brokers=localhost:9092
