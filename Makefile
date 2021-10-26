EXECUTABLE="todosServer"
.PHONY: openapi
openapi-codegen:
	oapi-codegen -generate types -o gen/todo_types.gen.go -package todos api/todo.yml
	oapi-codegen -generate server -o gen/todo_server.gen.go -package todos api/todo.yml

go-run:
	go run cmd/$(EXECUTABLE)/main.go

.PHONY: build
build:  ## build the API server binary
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/$(EXECUTABLE)

.PHONY: build-docker
build-docker: ## build the API server as a docker image
	docker build -f cmd/server/Dockerfile -t server .

.PHONY: clean
clean: ## remove temporary files
	rm -rf server coverage.out coverage-all.out

