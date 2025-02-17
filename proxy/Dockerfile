FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/proxy ./cmd/proxy

EXPOSE 8080

CMD ["/app/bin/proxy"]