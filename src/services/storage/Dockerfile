FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o storage .

FROM alpine:3.22

WORKDIR /app

RUN mkdir -p /var/log/data/parsed /var/log/data/storage

COPY --from=builder /app/storage .

CMD ["./storage"]
