package main

import (
	"builder/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	fmt.Println("Test m8")

	router := mux.NewRouter()

	router.HandleFunc("/new-build", handlers.Status)

	http.ListenAndServe(":8080", router)
}
