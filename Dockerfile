FROM golang:1.23.2 AS builder

RUN apt-get update && \
    apt-get install -y \
    build-essential \
    librdkafka-dev \
    git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

COPY . ./

RUN go build -o /usr/bin/application ./cmd/app

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y \
    librdkafka1 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/bin/application /usr/bin/application


EXPOSE 8000

CMD ["/usr/bin/application"]
