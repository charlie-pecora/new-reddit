package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/charlie-pecora/new-reddit/application"
	"github.com/charlie-pecora/new-reddit/authenticator"
)

func main() {
	auth, err := authenticator.New()
	if err != nil {
		log.Fatal(err)
	}
	router := application.New(auth)

	port := 3000
	log.Println("Starting server on port", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}

type IndexData struct {
	Name string
}

var indexTemplate = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, IndexData{"Charlie"})
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
