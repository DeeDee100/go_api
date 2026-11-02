#  docker exec -ti 7f256cdbb88a psql -U appuser -d appdb

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# rodar no container: make migrate-up DB_HOST=postgres
DB_HOST ?= localhost
POSTGRES_PORT ?= 5432

DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable


.PHONY: migrate-up migrate-down migrate-create migrate-version

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version


.PHONY: docker-logs docker-up docker-down db-shell
docker-logs:
	docker compose logs -f api

docker-up:
	docker compose up -d

docker-down:
	docker compose down

db-shell:
	@CONTAINER_ID=$$(docker ps -q -f "name=go_postgres"); \
	if [ -z "$$CONTAINER_ID" ]; then \
		echo "Oops Could not find 'go_postgres'. Run 'make docker-up'"; \
		exit 1; \
	fi; \
	echo "→ connecting to container $$CONTAINER_ID"; \
	docker exec -ti $$CONTAINER_ID psql -U appuser -d appdb