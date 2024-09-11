.PHONY: up-docker-db down-docker-db
up-docker-db:
	docker compose up -d
down-docker-db:
	docker compose rm -s test-db
