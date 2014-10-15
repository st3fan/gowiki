package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	DEFAULT_LISTEN_ADDRESS = "0.0.0.0"
	DEFAULT_LISTEN_PORT    = 8080
	DEFAULT_PREFIX         = "/"
)

func PageHandler(w http.ResponseWriter, r *http.Request) {
	input, err := ioutil.ReadFile("README.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output := blackfriday.MarkdownCommon(input)

	w.Header().Set("Content-Type", "text/html")
	w.Write(output)
}

func setupRouter(r *mux.Router) error {
	r.HandleFunc("/{pageName}", PageHandler).Methods("GET")
	return nil
}

func main() {
	router := mux.NewRouter()
	if err := setupRouter(router.PathPrefix(DEFAULT_PREFIX).Subrouter()); err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", DEFAULT_LISTEN_ADDRESS, DEFAULT_LISTEN_PORT)
	log.Printf("Starting wiki server on http://%s%s", addr, DEFAULT_PREFIX)
	http.Handle("/", router)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
