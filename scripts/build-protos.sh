#!/bin/sh

protoc -I ./api/proto/schemas/ \
    ./api/proto/schemas/*.proto \
    --go_out==grpc:./api/proto/proto_gencode \
    --go-grpc_out=./api/proto/proto_gencode