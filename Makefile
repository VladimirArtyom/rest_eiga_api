PATH_FILE = ./cmd/api/

POSTGRES_PASSWORD = 123456
POSTGRES_USER = postgres
DOCKER_CONTAINER_NAME = postgres-db
DB_PORT = 5432
DB_NAME = greenlight
DB_USERNAME = greenlight
DB_PASSWORD = 123456


create-network:
	docker network create eiga-go-network
delete-network:
	docker network rm eiga-go-network 
create-postgres:
	docker run --name $(DOCKER_CONTAINER_NAME) \
	--network eiga-go-network \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-p $(DB_PORT):5432 \
	-d postgres:latest

init-postgres:
	docker exec -i $(DOCKER_CONTAINER_NAME) psql -U $(POSTGRES_USER) < $(CURDIR)/db-scripts/init.sql
	docker exec -i $(DOCKER_CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(DB_NAME) < $(CURDIR)/db-scripts/grant-access-to-user.sql
	
start-postgres:
	docker start $(DOCKER_CONTAINER_NAME)

stop-postgres:
	docker stop $(DOCKER_CONTAINER_NAME)

remove-postgres:
	docker rm $(DOCKER_CONTAINER_NAME)

login-postgres:
	docker exec -it $(DOCKER_CONTAINER_NAME) \
	psql -U $(POSTGRES_USER) \

login-db:
	docker exec -it $(DOCKER_CONTAINER_NAME) \
	psql -U $(DB_USERNAME) \
	-d $(DB_NAME)

# migrate

migrate-up:
	docker run --rm \
	--network eiga-go-network \
	-v $(CURDIR)/migrations:/migrations \
	migrate/migrate:v4.14.1 \
	-path=/migrations -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DOCKER_CONTAINER_NAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# Please change the migration version, essayer à utiliser goto pour changing migrations version
migrate-goto:
	docker run --rm \
	--network eiga-go-network \
	-v $(CURDIR)/migrations:/migrations \
	migrate/migrate:v4.14.1 \
	-path=/migrations -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DOCKER_CONTAINER_NAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable" goto 2

# Please change the migration version
migrate-down:
	docker run --rm \
	--network eiga-go-network \
	-v $(CURDIR)/migrations:/migrations \
	migrate/migrate:v4.14.1 \
	-path=/migrations -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DOCKER_CONTAINER_NAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

migrate-version:
	docker run --rm \
	--network eiga-go-network \
	-v $(CURDIR)/migrations:/migrations \
	migrate/migrate:v4.14.1 \
	-path=/migrations -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DOCKER_CONTAINER_NAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable" version

migrate-create-movies-table_1:
	docker run --rm \
	--network eiga-go-network \
	-v $(CURDIR)/migrations:/migrations \
	migrate/migrate:v4.14.1 \
	create -seq -ext=.sql -dir=/migrations create_movies_table

migrate-create-movies-check-constraint_2:
	docker run --rm \
	--network eiga-go-network \
	-v $(CURDIR)/migrations:/migrations \
	migrate/migrate:v4.14.1 \
	create -seq -ext=.sql -dir=/migrations add_movies_check_constraints

init-db: create-network create-postgres 
delete-db: stop-postgres remove-postgres delete-network
populate-db: init-postgres 

run:
	go run $(PATH_FILE)
