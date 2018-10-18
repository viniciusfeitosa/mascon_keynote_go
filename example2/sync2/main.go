package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getLastName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var lastName string
	if vars["id"] == "1" {
		lastName = "Pacheco"
	} else {
		lastName = "Nadie"
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(lastName))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sync/{id:[0-9]+}", getLastName)

	http.ListenAndServe(":5000", r)
}
