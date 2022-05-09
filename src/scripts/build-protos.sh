#!/bin/sh

protoc -I ./src/protos/inputfile/ \
    ./src/protos/inputfile/*.proto \
    --go_out==grpc:./ \
    --go-grpc_out=.

protoc -I ./src/protos/onboarding/ \
    ./src/protos/onboarding/*.proto \
    --go_out==grpc:./ \
    --go-grpc_out=.