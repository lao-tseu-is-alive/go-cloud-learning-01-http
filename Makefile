#!make
include .env
export $(shell sed 's/=.*//' .env)

#the name of your API
APP=todos
EXECUTABLE=$(APP)Server

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

#.PHONY: db-docker-start
#db-docker-start: ## start the database server in docker container
#	@mkdir -p test/data/postgres
#	docker run --rm --name postgres -v $(shell pwd)/test/data:/testdata \
#		-v $(shell pwd)/test/data/postgres:/var/lib/postgresql/data \
#		-e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=$(APP) -e POSTGRES_DB=$(APP) -d -p 5433:5433 postgres


.PHONY: db-docker-start
## db-docker-start:	start docker postgres server, create app user&db  in a container named go-$(APP)-postgres
db-docker-start:
	@mkdir -p test/data/postgres
# docker container inspect call will set $? to 1 if container does not exist (cannot inspect) but
# to 0 if it does exist (this respects stopped containers). So the run command will just be called
# in case container does not exist as expected.
	docker container inspect go-$(APP)-postgres > /dev/null 2>&1 || \
	echo "  >  Starting your postgresql container..."; \
	docker run --name go-$(APP)-postgres \
	-e POSTGRES_USER=$(APP) -e POSTGRES_PASSWORD=$(APP) -e POSTGRES_DB=$(APP) -d postgres

#.PHONY: db-docker-is-ready
### db-docker-is-ready:	will wait until postgres is ready to be used
#db-docker-is-ready: db-docker-start
#	@echo "  >  Waiting for postgresql to be ready"
#	DOCKER_CONTAINER_NAME=go-$(APP)-postgres \
#	bash -c "until docker exec $(DOCKER_CONTAINER_NAME) pg_isready ; do sleep 5 ; done"
#docker exec -it go-$(APP)-postgres pg_isready
#timeout 90s bash -c "until docker exec $(DOCKER_CONTAINER_NAME) pg_isready ; do sleep 5 ; done"

.PHONY: db-docker-psql
## db-docker-psql:	open a psql session  on your app database  in docker container
db-docker-psql: db-docker-start
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


.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run in  "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

