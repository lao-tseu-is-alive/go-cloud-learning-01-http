.PHONY: openapi
openapi-codegen:
	oapi-codegen -generate types -o todo_types.gen.go -package main todo.yml
	oapi-codegen -generate server -o todo_server.gen.go -package main todo.yml

go-run:
	go run main.go todo_*.go
