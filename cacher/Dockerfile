FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/cacher ./cmd/cacher

EXPOSE 8080
CMD ["/app/bin/cacher"]
