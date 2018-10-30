package handlers

import (
	"builder/middleware"
	"builder/models"
	"encoding/json"
	"net/http"
)

// Status is the system's health check.
func Status(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}

func Home(w http.ResponseWriter, req *http.Request) {
	context, _ := req.Context().Value(middleware.CustomContextKey).(middleware.CustomContext)

	context.Db.CreateBuild(&models.Build{ID: "Another", Status: models.STARTED})

	w.Write([]byte("Done"))
}

func CreateBuild(w http.ResponseWriter, req *http.Request) {
	newModel := &models.Build{ID: "Another build", Status: models.STARTED}

	json.NewEncoder(w).Encode(newModel)
}

func GetBuild(w http.ResponseWriter, req *http.Request) {
	// v := mux.Vars(req)

	// fmt.Fprintf(w, "Category: %v\n", v["category"])

	mapD := map[string]int{"apple": 5, "lettuce": 7}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(mapD)
}
