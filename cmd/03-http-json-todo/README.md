## Let's make your first Go API with a classical TODO

In this next iteration of our hello web server we will introduce those features:

1. Refactor the code using dedicated struct type to store information related to all handlers of web server
2. Implement graceful shuts down the server without interrupting any active connections.
3. Implement a default handler to intercepts and handle all request to nonexistent endpoints.
4. Learn how to configure and use a logger


### Install the tools and dependency

First let's install [oapi-codegen](https://github.com/deepmap/oapi-codegen) an OpenAPI Client and Server Code Generator for Go.

From your terminal, go out of this project tree and run 

```bash
#go to your home directory first
cd ~
#download the oapi-codegen
go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
```
Now come back to your source code, where you made
Because I have **already done** at the beginning of the project a mod init

*go mod github.com/lao-tseu-is-alive/go-cloud-learning-01-http*

you have already a file called "mod.go" in the base directory of this project.

This file keep track of all your external Go dependencies. 


```bash
#go to your source code directory first (where the go.mod is) 
cd -
#install the required dependencies found in your go.mod files
go mod download
```


### What to remember :
- It is important to handle graceful shutdown for long queries to close cleanly

```go

```

### How to run it ?

to run the server and listen on port 3333  just type :
```bash
PORT=3333 go run main.go 


```


### More information :
- [OpenAPI-Specification](https://github.com/OAI/OpenAPI-Specification)
- [Golang Server code generation from OpenAPI Specification](https://medium.com/sellerapp/golang-server-code-generation-from-openapi-specification-5b4e5aa7cee)
- [go-swagger : Todo List Tutorial old OPENAPI 2.0](https://github.com/go-swagger/go-swagger/blob/master/docs/tutorial/todo-list.md)
- 