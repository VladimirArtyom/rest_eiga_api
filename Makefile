PATH_FILE = ./cmd/api/

POSTGRES_PASSWORD = 123456
POSTGRES_USER = postgres
DOCKER_CONTAINER_NAME = postgres-db
DB_PORT = 5432
DB_NAME = greenlight
DB_USERNAME = greenlight


start-postgres:
	docker run --name $(DOCKER_CONTAINER_NAME) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-p $(DB_PORT):5432 \
	-d postgres:latest

login-postgres:
	docker exec -it $(DOCKER_CONTAINER_NAME) \
	psql -U $(POSTGRES_USER) \

login-db:
	docker exec -it $(DOCKER_CONTAINER_NAME) \
	psql -U $(DB_USERNAME) \
	-d $(DB_NAME)

run:
	go run $(PATH_FILE)
