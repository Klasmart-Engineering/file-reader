#!/bin/sh

protoc -I ./src/protos \
    ./src/protos/*.proto \
    --go_out==grpc:./