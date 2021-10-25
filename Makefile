.PHONY: openapi
openapi-codegen:
	oapi-codegen -generate types -o gen/todo_types.gen.go -package todos api/todo.yml
	oapi-codegen -generate server -o gen/todo_server.gen.go -package todos api/todo.yml

go-run:
	go run cmd/todosServer/main.go gen/todo_*.go
