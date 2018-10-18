package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var sync2URL = os.Getenv("SYNC2_URL")
var sync3URL = os.Getenv("SYNC3_URL")

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

	lastName, err := getData(sync2URL, u.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	twitter, err := getData(sync3URL, u.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.LastName = lastName
	u.Twitter = twitter

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

func getData(baseURL, id string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/sync/%s", baseURL, id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	r := mux.NewRouter()
	r.Handle("/sync/{id:[0-9]+}", User{})

	http.ListenAndServe(":5000", r)
}
