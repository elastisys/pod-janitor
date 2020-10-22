FROM golang:alpine AS builder
WORKDIR /go/src/github.com/filetrust/pod-janitor
COPY . .
RUN cd cmd \
    && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o  pod-janitor .

FROM scratch
COPY --from=builder /go/src/github.com/filetrust/pod-janitor/cmd/pod-janitor /bin/pod-janitor

ENTRYPOINT ["/bin/pod-janitor"]