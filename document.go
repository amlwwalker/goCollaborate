package main 

import (
	"log"
	"path/filepath"
	"net/http"
	"encoding/json"
	"strings"
	"os"
)

type Response struct {
	Code int
	Description string
}

func walkpath(path string, f os.FileInfo, err error) error {
	log.Printf("%s with %d bytes\n", path,f.Size())
	return nil
}


func newDocument(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("name")
	//create a new document
	f, err := os.Create(location + "/" + filename + ".md")
	check(err)
	defer f.Close()
	log.Println("document created: " + filename)

	var response Response
	response.Code = 200
	response.Description = "document created: " + filename

	js, err := json.Marshal(response)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func listDocuments(w http.ResponseWriter, r *http.Request) {
	//list all the documents and return them as json
	var fileArray []string

	filepath.Walk(location, func(path string, f os.FileInfo, err error) error {

			if f.Mode().IsDir() && f.Name() == ".git" {
				return filepath.SkipDir
			} else {
				filename := strings.Replace(path, location, "", -1)
				if filename != "" {
					fileArray = append(fileArray, filename)	
				}
			}
		return nil
	});
	js, err := json.Marshal(fileArray)
	check(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}