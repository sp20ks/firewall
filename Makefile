BINARY=bin/firewall
ROOT_FOLDER=./cmd/firewall

.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY}-darwin ${ROOT_FOLDER}
	GOARCH=amd64 GOOS=linux go build -o ${BINARY}-linux ${ROOT_FOLDER}
	GOARCH=amd64 GOOS=windows go build -o ${BINARY}-windows ${ROOT_FOLDER}

.PHONY: run
run:
	go run ${ROOT_FOLDER}

.PHONY: clean
clean:
	go clean
	rm -r -f bin

.PHONY: run-containers
run-containers:
	docker run --rm -d -p 9001:80 --name server1 kennethreitz/httpbin
	docker run --rm -d -p 9002:80 --name server2 kennethreitz/httpbin
	docker run --rm -d -p 9003:80 --name server3 kennethreitz/httpbin

.PHONY: stop
stop:
	docker stop server1
	docker stop server2
	docker stop server3
