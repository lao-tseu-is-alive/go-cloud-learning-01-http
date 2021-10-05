## Your second version of a golang web server

In this next iteration of our tiny web server we will introduce three small but important features:

1. Define the listening port from an env variable PORT if it exists and is valid number
2. Check if  query parameter "username" exists and is not empty, then use it to make a warm and personal greeting.
3. Use Go html/template to send back to the client a valid html response. 


### Get the Listening Port from env variable PORT

It is very important to understand that you cannot hardcode a listening port 
in today's world of MicroServices. 

Buy the way in the ["The Twelve-Factor App"](https://12factor.net/), the third rule tells you:

- [III. Config : Store config in the environment](https://12factor.net/config)

*In the main function you can add the next lines to check the ENV for a variable PORT.
The line containing [os.LookupEnv("PORT")](https://pkg.go.dev/os#LookupEnv) does the job :*

```go
package main
import (
	"fmt"
	"log"
	"os"
	"strconv"
)
const defaultPort = 8080
func main() {
	listenAddr := fmt.Sprintf(":%v", defaultPort)
	val, exist := os.LookupEnv("PORT")
	if exist {
		port, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal("ERROR: PORT ENV should contain a valid integer value !")
		}
		listenAddr = fmt.Sprintf(":%v", port)
	}
	log.Printf(" ### Starting server... try navigating to http://localhost%v/hello to be greeted", listenAddr)
}
```

### What to remember :
- Read and understand [III. Config : Store config in the environment](https://12factor.net/config)
- Even with a great [net/http](https://pkg.go.dev/net/http) package some basic repetitive tasks can be tedious. That's when a framework comes to help.- 


### How to run it ?

to run the server and listen on port 3333  just type :
```bash
PORT=3333 go run main.go 
#in another terminal test it with curl or open a browser
curl http://localhost:8080/hello
curl  http://localhost:8080/hello?username=Rob%20Pike
#check what happens if you use another HTTP verb like POST
curl -XPOST  http://localhost:8080/hello?username=Rob%20Pike
curl -XPOST -d '{"username":"toto"}' http://localhost:8080/hello?username=Rob%20Pike
curl -XPUT -d '{"username":"toto"}' http://localhost:8080/1234
# what appears in the log for the last line ?
# we are using the default Mux of GO http, so for now we have no code to intercept this call

```

to build a binary executable called "mywebserver", just type :
```bash
go build -o mywebserver main.go 
```

to unit test the handler just type, by the way have a look at the main_test.go code :
```bash
go test -v 
```

to check on your server which process is listening on a specific port :
```bash
ss -lntp
```


### More information :
- [Demystifying HTTP Handlers in Golang](https://medium.com/geekculture/demystifying-http-handlers-in-golang-a363e4222756)