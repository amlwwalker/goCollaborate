package main

import (
	"log"
	"net/http"
	"fmt"
)

var location string
func main() {

	location = "/Users/alex/go/src/github.com/files/"
	userRepo.openRepository(location)
	fmt.Println("Starting in: " + location)

	go h.run()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/new", newDocument) //make a new document
	http.HandleFunc("/list", listDocuments) //make a new document
	log.Fatal(http.ListenAndServe(":8080", nil))
}
