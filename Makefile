PATH_FILE = ./cmd/api
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

migrate-force-cnum:
	docker run --rm \
		--network eiga-go-network \
		-v $(CURDIR)/migrations:/migrations \
		migrate/migrate:v4.14.1 \
		-path=/migrations -database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DOCKER_CONTAINER_NAME):$(DB_PORT)/$(DB_NAME)?sslmode=disable" force 4 # force your_num; Please change le numero 

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

migrate-add-movies-indexes_3:
	docker run --rm \
		--network eiga-go-network \
		-v $(CURDIR)/migrations:/migrations \
		migrate/migrate:v4.14.1 \
		create -seq -ext=.sql -dir=/migrations add_movies_indexes

migrate-create-users-table_4:
	docker run --rm \
		--network eiga-go-network \
		-v $(CURDIR)/migrations:/migrations \
		migrate/migrate:v4.14.1 create -seq -ext=.sql -dir=/migrations create_users_table

migrate-create-tokens-table_5:
	docker run --rm \
		--network eiga-go-network \
		-v $(CURDIR)/migrations:/migrations \
		migrate/migrate:v4.14.1 create -seq -ext=.sql -dir=/migrations create_tokens_table

init-db: create-network create-postgres 
delete-db: stop-postgres remove-postgres delete-network
populate-db: init-postgres 

run_concurrency:
	printf '%s\n' {1..8} | xargs -I % -P 8 curl -X PATCH -d '{"runtime":"78 mins"}' "http://localhost:8080/v1/movies/5" -printf '%s\n' {1..8} | xargs -I % -P 8 curl -X PATCH -d '{"runtime":"78 mins"}' "http://localhost:8080/v1/movies/5"

run:
	go run $(PATH_FILE)

## Requesting at the same times to test simultaneously req
same-occurance:
	curl localhost:8080/v1/movies/1 & curl localhost:8080/v1/movies/1 &

## needs to run independently,  このコードはcannot run using the make approach.
shot-api-simultaneously:
	for i in {1..6}; do curl http://localhost:8080/v1/healthcheck;done

## Authorization test　を 始める
get_bearer:
	curl -d '{"email": "konas@memes.com", "password":"password"}' localhost:8080/v1/tokens/authentication 
try_bearer:
	curl -i -H "Authorization: Bearer BWEP6MYF5ZPZNBBKCYPNUTY5KQ" localhost:8080/v1/healthcheck
try_invalid_bearer:
	curl -i -H "Authorization: Bearer MEMES" localhost:8080/v1/healthcheck
## Authorization Test を　終わる
