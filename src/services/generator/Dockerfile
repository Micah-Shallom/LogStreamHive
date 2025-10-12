FROM golang:1.24.5-alpine3.22

WORKDIR /app

#copy go.mod and go.sum files
COPY go.* ./   
RUN go mod download

COPY . .

#build
RUN go build -o /bin/logger . 

EXPOSE 8080

RUN mkdir -p /var/log/logger

CMD ["/bin/logger"]