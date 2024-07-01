BINARY=bin/firewall
ROOT_FOLDER=./cmd/firewall

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY}-darwin ${ROOT_FOLDER}
	GOARCH=amd64 GOOS=linux go build -o ${BINARY}-linux ${ROOT_FOLDER}
	GOARCH=amd64 GOOS=windows go build -o ${BINARY}-windows ${ROOT_FOLDER}
run:
	go run ${ROOT_FOLDER}
clean:
	go clean
	rm -r -f bin
