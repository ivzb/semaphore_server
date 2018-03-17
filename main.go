package main

import (
	"log"
	"strconv"

	"github.com/ivzb/achievers_server/app/controller"
	"github.com/ivzb/achievers_server/app/db"
	//"github.com/ivzb/semaphore_server/app/db"
	"github.com/ivzb/semaphore_server/app/middleware/app"
	"github.com/ivzb/semaphore_server/app/middleware/auth"
	"github.com/ivzb/semaphore_server/app/middleware/logger"
	"github.com/ivzb/semaphore_server/app/shared/config"
	"github.com/ivzb/semaphore_server/app/shared/file"
	l "github.com/ivzb/semaphore_server/app/shared/logger"
	"github.com/ivzb/semaphore_server/app/shared/token"
	"github.com/ivzb/semaphore_server/app/shared/uuid"

	"net/http"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	confBytes, err := file.Read("config" + string(os.PathSeparator) + "config.json")

	if err != nil {
		log.Panic(err)
	}

	conf, err := config.New(confBytes)

	if err != nil {
		log.Panic(err)
	}

	db, err := db.NewDB(conf.Database)

	if err != nil {
		log.Panic(err)
	}

	token, err := token.NewTokener(conf.Token)

	if err != nil {
		log.Panic(err)
	}

	logger := l.NewLogger()

	uuid := uuid.NewUUID()

	env := &a.Env{
		DB:     db,
		Log:    logger,
		Token:  token,
		Config: conf,
		UUID:   uuid,
	}

	http.Handle("/", anonChain(env, controller.HomeIndex))

	http.Handle("/"+conf.Server.Version+"/user/auth", anonChain(env, controller.UserAuth))
	http.Handle("/"+conf.Server.Version+"/user/create", anonChain(env, controller.UserCreate))

	port := strconv.Itoa(conf.Server.HTTPPort)
	logger.Message("started@:" + port)
	http.ListenAndServe(":"+port, nil)
}

func authChain(env *a.Env, handler app.Handler) http.Handler {
	return use(app.App{env, handler}, auth.Handler, logger.Handler)
}

func anonChain(env *a.Env, handler app.Handler) http.Handler {
	return use(app.App{env, handler}, logger.Handler)
}

// specify middlewares in reverse order since it is chaining them recursively
func use(app app.App, middlewares ...func(app.App) app.App) http.Handler {
	for _, middleware := range middlewares {
		app = middleware(app)
	}

	var handler http.Handler = app

	return handler
}
