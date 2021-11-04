# go-cloud-learning-01-http
From zero to hero : Learning how to build your first «Todos» api server in Go.



## Requirements
 1. Have access to a Linux Ubuntu 20.04 box or VM 
 2. Install GO :  https://golang.org/doc/install
 3. Pick your favorite editor like vim or maybe [vscode](https://code.visualstudio.com/docs/setup/linux)
 4. Get a basic knowledge of GO : https://go.dev/learn/


## Project Layout and conventions
This project uses the Standard Go Project Layout : https://github.com/golang-standards/project-layout

## Where do I start ?
1. just clone this repo in you Linux Box : git clone https://github.com/lao-tseu-is-alive/go-cloud-learning-01-http.git
2. cd into the cloned repo directory
3. install all project dependencies with :  **go get -d ./...**
4. then jump directly to your first example in [cmd/00-http-basic](https://github.com/lao-tseu-is-alive/go-cloud-learning-01-http/tree/main/cmd/00-http-basic)
5. have a look to the code  in main.go file
6. just run it : **go run main.go**
7. hack-it, modify-it and have fun !
8. Follow along with the other 3 examples in cmd directory:
   1. [01-http-envconfig-template-parameter](https://github.com/lao-tseu-is-alive/go-cloud-learning-01-http/tree/main/cmd/01-http-envconfig-template-parameter) : *learn how to use env variables for config and reading a parameter*
   2. [02-http-refactor-graceful-shutdown](https://github.com/lao-tseu-is-alive/go-cloud-learning-01-http/tree/main/cmd/02-http-refactor-graceful-shutdown) : *learn about using a logger and code like a pro implementing «graceful» shutdown* 
   3. [03-http-json-todo](https://github.com/lao-tseu-is-alive/go-cloud-learning-01-http/tree/main/cmd/03-http-json-todo) : *implement a basic "todos" API using [Echo](https://echo.labstack.com/) micro framework*.
9.When you are ready, jump to the main example of this TODOS Api Server

## Main example is a template Go project
You can use this repository as a base template for your future projects. 
You can try it with :

```bash
make db-docker-init-data
make run
``` 
The first command will start a docker postgres container, create a database, 
do a db migration up to create the necessary tables, and load some test data.

The next command will compile and run your todosServer, injecting version info based on your git tag, 
and using your .env files to initialize the env variables.  

**The main features of this template are :**
+ A Makefile with more than 14 ready to use sub-commands _(you can try : **make help**)_.    
+ [Echo](https://echo.labstack.com/) : *High performance, extensible, minimalist Go web framework*
+ Contract based development using [OpenAPI 3](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md), *(also called swagger in the past)*.
+ Go server code generation from OpenApi spec using [oapi-codegen](https://github.com/deepmap/oapi-codegen)
+ Database migration using [golang-migrate](https://github.com/golang-migrate/migrate)
+ Table Driven Testing 
+ Live reload  with **make   [fswatch](https://github.com/emcrisostomo/fswatch)
+ Server version defined automatically based on your git tags [semantic versioning](https://semver.org/). For example 0.1.1  **git tag -a v0.1.1 -m "v0.1.1"**  

## Useful Links
- [Golang.org : home for the official project](https://golang.org/)
- [A tour of Go is an interactive tour of Go ](https://tour.golang.org/)
- [Go.dev a hub for Go developers](https://go.dev/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)
- [Go by Example](https://gobyexample.com/)
- [Golang tutorial and examples](https://www.golangprograms.com/golang-package-examples.html)
- [My golang learning](https://github.com/lao-tseu-is-alive/golang-learning)
- [Why you should try Go as a PHP Developer(YouTube video)](https://www.youtube.com/watch?v=Mjcw8fHdx8Q)
- [How I write HTTP services after eight years by Mat Ryer](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html)
- [Alternate OpenAPI Generator](https://github.com/OpenAPITools/openapi-generator) 
