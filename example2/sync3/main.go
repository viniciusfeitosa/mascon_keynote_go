package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func getTwitter(w http.ResponseWriter, r *http.Request) {
	time.Sleep(50 * time.Millisecond)
	vars := mux.Vars(r)
	var twitter string
	if vars["id"] == "1" {
		twitter = "@viniciuspach"
	} else {
		twitter = "@nose"
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(twitter))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/sync/{id:[0-9]+}", getTwitter)

	http.ListenAndServe(":5000", r)
}
