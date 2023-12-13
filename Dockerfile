FROM golang:1.21

WORKDIR /app
COPY . .
RUN go mod download
EXPOSE 1234
RUN go run main.go
