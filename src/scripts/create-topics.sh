#!/bin/sh

<<<<<<< HEAD
docker exec -it redpanda-1 rpk topic create organization --brokers=localhost:9092
docker exec -it redpanda-1 rpk topic create S3FileCreatedUpdated --brokers=localhost:9092
docker exec -it redpanda-1 rpk topic create organization_proto --brokers=localhost:9092
docker exec -it redpanda-1 rpk topic create organization-proto-12 --brokers=localhost:9092
=======
docker exec -it redpanda-1 rpk topic create organization_proto --brokers=localhost:9092
>>>>>>> 7267a7f (Add proto schema cache, refactor proto service)
