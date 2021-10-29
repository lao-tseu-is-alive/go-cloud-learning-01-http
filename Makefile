#!make
include .env
# next line allows to export env variables to external process (like your Go app)
export $(shell sed 's/=.*//' .env)

#the name of your API
APP=todos
EXECUTABLE=$(APP)Server
APP_DSN=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
# using golang-migrate https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
# here with the docker file so no need to install it
# MIGRATE := docker run -v $(shell pwd)/db/migrations:/migrations --network host migrate/migrate:v4.10.0 -path=/db/migrations/ -database "$(APP_DSN)"
# or download your release from here : https://github.com/golang-migrate/migrate/releases
# for ubuntu & debian : wget https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.deb
MIGRATE=/usr/local/bin/migrate -database "$(APP_DSN)" -path=db/migrations/

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent


# for reason to use .Phony see : https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: openapi
## openapi-codegen:	will generate helper Go code for types & server based on OpenApi spec in api/app.yml
openapi-codegen:
	oapi-codegen -generate types -o internal/todos/todo_types.gen.go -package todos api/todos.yml
	oapi-codegen -generate server -o internal/todos/todo_server.gen.go -package todos api/todos.yml

.PHONY: run
## build:	will compile your server app binary and place it in the bin sub-folder
run:
	go run cmd/$(EXECUTABLE)/main.go

.PHONY: build
## build:	will compile your server app binary and place it in the bin sub-folder
build:
	@echo "  >  Building your app binary inside bin directory..."
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o bin/$(EXECUTABLE) cmd/$(EXECUTABLE)/main.go

.PHONY: build-docker
build-docker: ## build the API server as a docker image
	docker build -f cmd/server/Dockerfile -t server .

.PHONY: clean
## clean:	will delete you server app binary and remove temporary files like coverage output
clean:
	rm -rf bin/$(EXECUTABLE) coverage.out coverage-all.out

.PHONY: db-docker-start
## db-docker-start:	start docker postgres server, create app user&db  in a container named go-$(APP)-postgres
db-docker-start:
	@mkdir -p test/data/postgres
# docker container inspect call will set $? to 1 if container does not exist (cannot inspect) but
# to 0 if it does exist (this respects stopped containers). So the run command will just be called
# in case container does not exist as expected.
# -v $(shell pwd)/testdata/postgres:/var/lib/postgresql/data
	@echo "  >  Checking if for postgresql container is already there"
	docker container inspect go-$(APP)-postgres > /dev/null 2>&1 || \
	(echo "  >  Started your postgresql container on port ${DB_HOST}:${DB_PORT}..."; \
	docker run --name go-$(APP)-postgres \
	-v $(shell pwd)/test/data:/testdata  \
	-e POSTGRES_USER=$(APP) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(APP) -d -p $(DB_HOST):$(DB_PORT):5432 postgres \
	)


.PHONY: db-docker-is-ready
## db-docker-is-ready:	will wait until postgres is ready to be used
db-docker-is-ready: db-docker-start
	@echo "  >  Waiting for postgresql to be ready"
	until docker exec go-$(APP)-postgres pg_isready ; do sleep 3 ; done
#docker exec -it go-$(APP)-postgres pg_isready
#timeout 90s bash -c "until docker exec $(DOCKER_CONTAINER_NAME) pg_isready ; do sleep 5 ; done"

.PHONY: db-docker-psql
## db-docker-psql:	open a psql session  on your app database  (starts docker container if not running)
db-docker-psql: db-docker-start db-docker-is-ready
	@echo "  >  You can also use : scripts/connect_psql_with_dot_env.sh"
	@echo "  >  Entering in postgresql psql shell (use \q to exit) "
	docker exec -it go-$(APP)-postgres psql $(APP) -U $(APP)

.PHONY: db-docker-stop
## db-docker-stop:	stop the database server in docker container
db-docker-stop:
	@echo "  >  Stopping your postgresql container... "
	docker stop go-$(APP)-postgres

.PHONY: db-docker-delete
## db-docker-delete:	clean the database directory and removes docker container
db-docker-delete: db-docker-stop
	@echo "  >  Removing your postgresql container... "
	@rm -rf test/data/postgres/*
	@rmdir test/data/postgres
	docker rm go-$(APP)-postgres

.PHONY: db-docker-migrate-new
## db-docker-migrate-new:	create a new database migration
db-docker-migrate-new:
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir db/migrations/ $${name// /_}

.PHONY: db-docker-migrate-up
## db-docker-migrate-up: run all new database migrations
db-docker-migrate-up:
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: db-docker-migrate-down
## db-docker-migrate-up: revert database to the last migration step
db-docker-migrate-down:
	@echo "Running all new database migrations..."
	@$(MIGRATE) down 1

.PHONY: db-docker-migrate-reset
## db-docker-migrate-reset:	 reset database and re-run all migrations
db-docker-migrate-reset:
	@echo "Resetting database..."
	@$(MIGRATE) drop
	@echo "Running all database migrations..."
	@$(MIGRATE) up


.PHONY: db-docker-init-data
## db-docker-init-data:	load initial data in your app database  (starts docker container if not running)
db-docker-init-data: db-docker-start db-docker-is-ready
	@echo "  >  Loading data in your postgresql db... "
	docker exec -it go-$(APP)-postgres psql $(APP) -U $(APP) -f /testdata/initial_todos_data.sql



.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run in  "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

