NETWORK_NAME=firewall-network
PROXY_IMAGE=proxy
PROXY_CONTAINER=proxy
PROXY_PORT=8080

.PHONY: build
build: build-proxy

.PHONY: build-proxy
build-proxy:
	docker build -t $(PROXY_IMAGE) ./proxy

.PHONY: network
network:
	@docker network ls | grep -q $(NETWORK_NAME) || docker network create $(NETWORK_NAME)

.PHONY: run-containers
run-containers: network
	docker run --rm -d --network $(NETWORK_NAME) --name server1 -p 9001:80 kennethreitz/httpbin
	docker run --rm -d --network $(NETWORK_NAME) --name server2 -p 9002:80 kennethreitz/httpbin
	docker run --rm -d --network $(NETWORK_NAME) --name server3 -p 9003:80 kennethreitz/httpbin
	docker run --rm -d --network $(NETWORK_NAME) -p 8080:8080 --name $(PROXY_CONTAINER) $(PROXY_IMAGE)

.PHONY: stop
stop:
	docker stop server1 server2 server3 $(PROXY_CONTAINER) || true

.PHONY: clean
clean: stop
	@docker network rm $(NETWORK_NAME) || true
	@docker rmi $(PROXY_IMAGE) || true

.PHONY: logs-proxy
logs-proxy:
	docker logs -f $(PROXY_CONTAINER)
