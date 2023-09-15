FROM golang:1.17

WORKDIR /app

COPY kafka_listener/go.mod kafka_listener/go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y build-essential git librdkafka-dev

COPY kafka_listener/ ./
RUN go build -o app

CMD ["./app"]