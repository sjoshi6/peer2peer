# golang image where workspace (GOPATH) configured at /go.
FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/sjoshi6/peer2peer

# Build the golang-docker command inside the container.
RUN go install github.com/sjoshi6/peer2peer

# Run the golang-docker command when the container starts.
ENTRYPOINT ["/go/bin/peer2peer api"]

# http server listens on port 8000.
EXPOSE 8000
