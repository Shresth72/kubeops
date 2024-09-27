FROM golang:1.20 AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

RUN go build -o /kluster .

FROM alpine:latest

COPY --from=builder /kluster /usr/local/bin/kluster 

ENTRYPOINT ["/usr/local/bin/kluster"]
