FROM golang:1.24.5-alpine3.22  AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o log_parser .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/log_parser .

RUN mkdir -p /logs /parsed_logs

ENV LOG_DIR=/logs
ENV PARSED_LOG_DIR=/parsed_logs

RUN chmod +x /app/log_parser

CMD ["./log_parser"]
