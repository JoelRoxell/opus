package server

import (
	"builder/handlers"
	"builder/middleware"
	"builder/repositories"
	"log"
	"net/http"
	"os"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server will provide same availability as the CLI, and will be able to listen for remote hooks.
func Server(exit chan bool) {
	// TODO: Allow same actions as via CLI, this section is not used atm...
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.Home)
	router.HandleFunc("/health", handlers.Status)
	router.HandleFunc("/build", handlers.CreateBuild).Methods("POST")
	router.HandleFunc("/build/{id}", handlers.GetBuild).Methods("GET")

	loggerRouter := ghandlers.LoggingHandler(os.Stdout, router)

	session, _ := repositories.NewMongoDBConnection("localhost")

	customContxt := &middleware.CustomContext{Db: session}

	srv := &http.Server{Addr: ":8888", Handler: middleware.ContextHandler(loggerRouter.ServeHTTP, *customContxt)}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	<-exit

	srv.Shutdown(nil)
}
