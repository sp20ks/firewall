NETWORK_NAME=firewall-network
FIREWALL_IMAGE=firewall
FIREWALL_CONTAINER=firewall

.PHONY: build
build:
	docker build -t $(FIREWALL_IMAGE) .

.PHONY: network
network:
	@docker network ls | grep -q $(NETWORK_NAME) || docker network create $(NETWORK_NAME)

.PHONY: run-containers
run-containers: network
	docker run --rm -d --network $(NETWORK_NAME) --name server1 -p 9001:80 kennethreitz/httpbin
	docker run --rm -d --network $(NETWORK_NAME) --name server2 -p 9002:80 kennethreitz/httpbin
	docker run --rm -d --network $(NETWORK_NAME) --name server3 -p 9003:80 kennethreitz/httpbin
	docker run --rm -d --network $(NETWORK_NAME) -p 8080:8080 --name $(FIREWALL_CONTAINER) $(FIREWALL_IMAGE)

.PHONY: stop
stop:
	docker stop server1 server2 server3 $(FIREWALL_CONTAINER) || true

.PHONY: clean
clean: stop
	docker network rm $(NETWORK_NAME) || true
	docker rmi $(FIREWALL_IMAGE) || true

.PHONY: logs
logs:
	docker logs -f $(FIREWALL_CONTAINER)
