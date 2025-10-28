FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o tcp_server .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/tcp_server .

RUN chmod +x /app/tcp_server

EXPOSE 9000

CMD ["./tcp_server"]



