FROM golang:1.24.5-alpine3.22  AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o log_query .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/log_query .

RUN chmod +x /app/log_query

CMD ["./log_query"]
