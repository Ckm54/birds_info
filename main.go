package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Bird struct {
	Species			string	`json:"species"`
	Description	string	`json:"description"`
}

var birds []Bird

// creates router and returns so we can test the
// router outside the main function
func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")
	r.HandleFunc("/bird", getBirdHandler).Methods("GET")
	r.HandleFunc("/bird", createBirdHandler).Methods("POST")

	staticFileDirectory := http.Dir("./assets")

	staticFileHandler := http.StripPrefix("/assets/" ,http.FileServer(staticFileDirectory))

	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	return r
}

func main() {
	r := newRouter()

	http.ListenAndServe(":8080", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func getBirdHandler(w http.ResponseWriter, r *http.Request) {
	birdListBytes, err := json.Marshal(birds)

	// if error print to console and return server error to the user
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// else write json list of birds to the response
	w.Write(birdListBytes)
}

func createBirdHandler(w http.ResponseWriter, r *http.Request) {
	bird := Bird{}

	// parse the form values
	err := r.ParseForm()

	// handle error
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bird.Species = r.Form.Get("species")
	bird.Description = r.Form.Get("description")

	// append existing list of birds with a new entry
	birds = append(birds, bird)

	// redirect user to the original page located at /assets/
	http.Redirect(w, r, "/assets/", http.StatusFound)
}