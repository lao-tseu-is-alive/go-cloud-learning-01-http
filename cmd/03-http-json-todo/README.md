## Let's make your first Go API with a classical TODO

Congratulations, now you will learn how to write your first MicroServices. 
For education purposes it will be a basic TODO app with data stored in memory.
But in the process you will grasp some new concepts


 1. From now on, we will start using [Echo](https://echo.labstack.com/) as web framework.
 2. We will use [Go Modules](https://github.com/golang/go/wiki/Modules) dependency manager
 3. We will introduce the concept of Contract based API using [OpenApi 3.0](https://github.com/OAI/OpenAPI-Specification) specification
 4. We will see  the [online swagger editor](https://editor.swagger.io/) to create/edit/validate OpenAPI 3 specifications
 5. We will generate code based on the OpenAPi specification using [oapi-codegen](https://github.com/deepmap/oapi-codegen)

### Install the tools and dependency

First let's install [oapi-codegen](https://github.com/deepmap/oapi-codegen) an OpenAPI Client and Server Code Generator for Go.

From your terminal, go out of this project tree and run 

```bash
#go to your home directory first
cd ~
#download the oapi-codegen
go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
```
Now come back to your source code tree.
Because I have **already done** at the beginning of this project a mod init :

*go mod github.com/lao-tseu-is-alive/go-cloud-learning-01-http*

you don't need to run this command in this project
you have already a file called "mod.go" in the base directory of this project.

This file keep track of all your external Go dependencies. 
Just run the next command to install the dependencies locally.
```bash
#go to your source code directory first (where the go.mod is) 
cd -
#install the required dependencies found in your go.mod files
go mod download
```


### What to remember :
- Design your API specification first and have the spec.yaml or spec.json in your code repository
- Try to follow the [REST API Guidelines](https://github.com/microsoft/api-guidelines/blob/vNext/Guidelines.md) as much as possible
- Be consistent in all your API definition.
- Generate the Go server and types from your OpenApi spec, use them as is, **DO NOT EDIT THEM !**
- Use low revision number of your API if you extend or add new path. 
- Use New major version if you introduce breaking changes in your API (remove or change something)
- 

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
- [Swagger Editor](https://editor.swagger.io/)
- 