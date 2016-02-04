package main

import (
	"log"
	"net/http"
	"fmt"
	// "github.com/libgit2/git2go"
)

// var userRepo Repo

func main() {
	fmt.Println("whats this")
	

	//lets open the repository
	// userRepo.repo = openRepository()

	go h.run()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
