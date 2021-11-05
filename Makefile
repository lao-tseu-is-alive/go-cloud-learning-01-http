#!make
ifneq ("$(wildcard $(.env))","")
	ENV_EXISTS = "TRUE"
	include .env
	# next line allows to export env variables to external process (like your Go app)
	export $(shell sed 's/=.*//' .env)
else
	ENV_EXISTS = "FALSE"
	DB_DRIVER=postgres
	DB_HOST=127.0.0.1
	DB_PORT=5432
	DB_NAME=todos
	DB_USER=todos
	# DB_PASSWORD should be defined in your env or in github secrets
	DB_SSL_MODE=disable
endif

#the name of your API
APP=todos
EXECUTABLE=$(APP)Server
# The version is based on your git tags, when you're ready make a new tag and push it to github et voila !
# git tag -a v0.1.0 -m "v0.1.0"
# git push origin main --tags
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "0.0.0_alpha")
REVISION ?= $(shell git rev-list -1 HEAD)
BUILD ?= $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
PACKAGES := $(shell go list ./... | grep -v /vendor/)
LDFLAGS := -ldflags "-X main.VERSION=${VERSION} -X main.GitRevision=${REVISION} -X main.BuildStamp=${BUILD}"
PID_FILE := "./$(APP).pid"
APP_DSN=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
ifeq ($(ENV_EXISTS),"TRUE")
	# or download your release from here : https://github.com/golang-migrate/migrate/releases
	# for ubuntu & debian : wget https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.deb
	MIGRATE=/usr/local/bin/migrate -database "$(APP_DSN)" -path=db/migrations/
else
	# using golang-migrate https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
	# here with the docker file so no need to install it
	MIGRATE := docker run -v $(shell pwd)/db/migrations:/migrations --network host migrate/migrate:v4.10.0 -path=/db/migrations/ -database "$(APP_DSN)"
endif

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# check we have a couple of dependencies
.PHONY: dependencies-openapi
dependencies-openapi:
	@command -v oapi-codegen >/dev/null 2>&1 || { printf >&2 "oapi-codegen is not installed, please run: go get github.com/jteeuwen/go-bindata/...\n"; exit 1; }

.PHONY: dependencies-fswatch
dependencies-fswatch:
	@command -v fswatch --version >/dev/null 2>&1 || { printf >&2 "fswatch is not installed, install it to be able to use 'make run-hot-reload'\n"; exit 1; }



# for reason to use .Phony see : https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: openapi-codegen
## openapi-codegen:	will generate helper Go code for types & server based on OpenApi spec in api/app.yml
openapi-codegen: dependencies-openapi
	oapi-codegen -generate types -o internal/todos/todo_types.gen.go -package todos api/todos.yml
	oapi-codegen -generate server -o internal/todos/todo_server.gen.go -package todos api/todos.yml

.PHONY: lint
## lint:	run golint on all your Go package
lint:
	@golint $(PACKAGES)


.PHONY: test
## test:	will run unit tests for all your Go code
test:
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test-cover
## test-cover:	will run all unit tests and show test coverage information
test-cover: test
	go tool cover -html=coverage-all.out

.PHONY: run
## run:	will run a dev version of your Go application
run:
	go run ${LDFLAGS} cmd/$(EXECUTABLE)/main.go

.PHONY: run-restart
# run-restart:	DO NOT USE will restart your server USE instead  : make run-live-reload
run-restart:
	@pkill -P `cat $(PID_FILE)` || true
	@printf '%*s\n' "80" '' | tr ' ' -
	@echo "Your code has changed. Will restart the server..."
	@go run ${LDFLAGS} cmd/$(EXECUTABLE)/main.go & echo $$! > $(PID_FILE)
	@printf '%*s\n' "80" '' | tr ' ' -


# you can find instructions on fswatch here : https://github.com/emcrisostomo/fswatch
.PHONY: run-hot-reload
## run-hot-reload: 	will run a dev version of your Go application with «live» ️reload 👍 😃 😋 (requires fswatch on your box)
run-hot-reload: dependencies-fswatch
	@go run ${LDFLAGS} cmd/$(EXECUTABLE)/main.go & echo $$! > $(PID_FILE)
	@fswatch -x -o --event Created --event Updated --event Renamed -r internal pkg cmd config | xargs -n1 -I {} make run-restart

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
	-e POSTGRES_USER=$(APP) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(APP) -d -p $(DB_HOST):$(DB_PORT):5432 postgres:14.0 \
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
## db-docker-migrate-new:	create a new database migration using https://github.com/golang-migrate/migrate
db-docker-migrate-new:
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir db/migrations/ $${name// /_}

.PHONY: db-docker-migrate-up
## db-docker-migrate-up: 	run all new database migrations
db-docker-migrate-up: db-docker-start db-docker-is-ready
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: db-docker-migrate-down
## db-docker-migrate-down: 	revert database to the last migration step
db-docker-migrate-down:
	@echo "Running all new database migrations..."
	@$(MIGRATE) down 1

.PHONY: db-docker-migrate-reset
## db-docker-migrate-reset:	reset database and re-run all migrations
db-docker-migrate-reset:
	@echo "Resetting database..."
	@$(MIGRATE) drop
	@echo "Running all database migrations..."
	@$(MIGRATE) up


.PHONY: db-docker-init-data
## db-docker-init-data:	load initial data in your app database  (starts docker container if not running)
db-docker-init-data: db-docker-migrate-up
	@echo "  >  Loading data in your postgresql db... "
	docker exec -it go-$(APP)-postgres psql $(APP) -U $(APP) -f /testdata/initial_todos_data.sql

.PHONY: db-github-migrate-up
## db-github-migrate-up: 	to be used in github actions to run all new database migrations
db-github-migrate-up:
	@echo "Running your database migrations..."
	@$(MIGRATE) up


.PHONY: db-github-init-data
## db-github-init-data:	to be used in github actions to load initial data in your app database
db-github-init-data:
	@echo "  >  Loading initial data in your postgresql db... "
	docker exec -it postgres psql $(APP) -U $(APP) -f /testdata/initial_todos_data.sql


.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run in  "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

