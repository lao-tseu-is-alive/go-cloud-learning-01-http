package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/internal/todos"
	"github.com/lao-tseu-is-alive/go-cloud-learning-01-http/pkg/config"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	appName           = "todos"
	defaultServerPort = 8080
	defaultServerIp   = "127.0.0.1"
	defaultDBPort     = 5432
	defaultDBIp       = "127.0.0.1"
	defaultDBPassword = "todos_password"
	/*  defaultDBSslMode: in dev env it can be ok to disable SSL mode but in prod it is another story
			it depends on various factor. is your service (go) running in the same host as the db (localhost ?)
	 		if not, is the network between your server and your db trusted ?? read the doc and ask your security officer:
	 		https://www.postgresql.org/docs/11/libpq-ssl.html#LIBPQ-SSL-PROTECTION
	*/
	defaultDBSslMode = "disable"
	defaultDBDriver  = "postgres"
	webRootDir       = "swagger-ui"
	/*
		shutDownTimeout     = 2 * time.Second // number of second to wait before closing server
		defaultReadTimeout  = 2 * time.Minute
		defaultWriteTimeout = 2 * time.Minute
		defaultWebRootDir   = "./web/dist"
	*/
)

// VERSION  the current version of the application is evaluated at build time based on your git tag.
var VERSION = "0.0.0"
var GitRevision = "unknown"
var BuildStamp = "unknown"

/*
//go:embed ./web/dist

*/
var embededFiles embed.FS

func getFileSystem(useOS bool, log *log.Logger) http.FileSystem {
	if useOS {
		log.Println("using live mode")
		return http.FS(os.DirFS("dist"))
	}

	log.Println("using embed mode")
	fsys, err := fs.Sub(embededFiles, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// GetNewServer initialize a new Echo server and returns it
func GetNewServer(useOS bool, l *log.Logger, store todos.Storage) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	myTodosApi := todos.Service{
		Log:   l,
		Store: store,
	}
	if useOS {
		webRootDirPath, err := filepath.Abs(webRootDir)
		if err != nil {
			log.Fatalf("Problem getting absolute path of directory: %s\nError:\n%v\n", webRootDir, err)
		}
		if _, err := os.Stat(webRootDirPath); os.IsNotExist(err) {
			log.Fatalf("The webRootDir parameter is wrong, %s is not a valid directory\nError:\n%v\n", webRootDirPath, err)
		}
		l.Printf("Using live mode serving from %s", webRootDirPath)
		e.Static("/", webRootDirPath)
	} else {
		assetHandler := http.FileServer(http.FS(embededFiles))
		e.GET("/", echo.WrapHandler(assetHandler))
		//e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
	}

	// here the routes defined in OpenApi todos.yaml are registered
	todos.RegisterHandlers(e, &myTodosApi)
	// add another route for maxId
	e.GET("/todos/maxid", myTodosApi.GetMaxId)
	return e
}

func GetVersion() string {
	return fmt.Sprintf("%s Ver: %s, Build: %s, rev: %s ", appName, VERSION, BuildStamp, GitRevision)
}

// main is the entry point of your todos Api TodosService service
func main() {
	displayVersion := flag.Bool("version", false, "display version and terminated")
	if *displayVersion {
		fmt.Printf("%s\n", GetVersion())
		os.Exit(0)
	} else {
		fmt.Printf("## Starting Go %s \n", GetVersion())
	}

	//l := log.New(ioutil.Discard, appName, 0)
	l := log.New(os.Stdout, appName, log.Ldate|log.Ltime|log.Lshortfile)

	//useOS := len(os.Args) > 1 && os.Args[1] == "live"
	// we cannot use embeded for now because it is not working yet
	useOS := true

	listenAddress, err := config.GetListenAddrFromEnv(defaultServerIp, defaultServerPort)
	if err != nil {
		log.Fatalf("💥💥 error doing config.GetListenAddrFromEnv. error: %v\n", err)
	}

	driver := config.GetDbDriverFromEnv(defaultDBDriver)
	if !todos.IsDriverSupported(driver) {
		log.Fatalf("💥💥 error the driver : %s is not supported yet.\n", driver)
	}

	var dbDsn = ""
	if driver == "postgres" {
		dbDsn, err = config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
			appName, appName, defaultDBPassword, defaultDBSslMode)
		if err != nil {
			log.Fatalf("💥💥 error doing config.GetPgDbDsnUrlFromEnv. error: %v\n", err)
		}
	}

	s, err := todos.GetStorageInstance(driver, dbDsn, l)
	if err != nil {
		l.Fatalf("💥💥 error getting Storage Instance for driver %s. error: %v\n", driver, err)
	}
	defer s.Close()

	e := GetNewServer(useOS, l, s)
	l.Printf("Will start http server ««%s»», listening on: %s \n", GetVersion(), listenAddress)
	e.Logger.Fatal(e.Start(listenAddress))
}
