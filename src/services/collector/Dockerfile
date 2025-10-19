FROM golang:1.24.5-alpine3.22 

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /bin/collector .

EXPOSE 8080

RUN mkdir -p /var/log/collector

RUN adduser -D -u 1000 collector
USER collector

# Command
CMD ["/bin/collector"]
