.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: logs-proxy
logs-proxy:
	docker-compose logs -f proxy

.PHONY: logs-ratelimiter
logs-ratelimiter:
	docker-compose logs -f ratelimiter

.PHONY: rebuild
rebuild:
	docker-compose up --build -d
