FROM golang:1.18-alpine

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/

COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN go build -ldflags "-s -w" -o main src/server/main.go

EXPOSE 6000
EXPOSE 9000

# Run the executable
CMD ["./main"]

