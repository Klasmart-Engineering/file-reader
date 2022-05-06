#!/bin/sh

protoc -I ./src/protos/csvfile/ \
    ./src/protos/csvfile/*.proto \
    --go_out==grpc:./ \
    --go-grpc_out=.

protoc -I ./src/protos/onboarding/ \
    ./src/protos/onboarding/*.proto \
    --go_out==grpc:./ \
    --go-grpc_out=.