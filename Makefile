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

.PHONY: logs-cacher
logs-cacher:
	docker-compose logs -f cacher

.PHONY: logs-auth
logs-auth:
	docker-compose logs -f auth

.PHONY: logs-rules-engine
logs-rules-engine:
	docker-compose logs -f rules-engine

.PHONY: rebuild
rebuild:
	docker-compose up --build -d
