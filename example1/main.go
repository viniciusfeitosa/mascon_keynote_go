package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// User is an the domain struct
type User struct {
	ID        string
	FirstName string
	LastName  string
	Twitter   string
}

func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u.ID = vars["id"]
	u.FirstName = domain1(u.ID)
	u.LastName = domain2(u.ID)
	u.Twitter = domain3(u.ID)

	data, err := json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func domain1(id string) string {
	if id == "1" {
		return "Vinicius"
	}
	return "Juan"
}

func domain2(id string) string {
	if id == "1" {
		return "Pacheco"
	}
	return "Nadie"
}

func domain3(id string) string {
	time.Sleep(50 * time.Millisecond)
	if id == "1" {
		return "@viniciuspach"
	}
	return "@nose"
}

func main() {
	r := mux.NewRouter()
	r.Handle("/rock/{id:[0-9]+}", User{})

	http.ListenAndServe(":5000", r)
}
