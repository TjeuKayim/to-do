package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"errors"
	"strings"
	"math/rand"
)

const urlPrefix = "/api"

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc(urlPrefix + "/page/", pageHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		//http.Error(w, "Test", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := fmt.Fprintln(w, "Hello World")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, urlPrefix)
	trimmedPath = strings.Trim(trimmedPath, "/")
	path := strings.Split(trimmedPath, "/")
	switch len(path) {
	case 1:
		if r.Method == "POST" {
			postPage(w, r)
		}
	case 2:
		id := PageId(path[1])
		switch r.Method {
		case "GET":
			getPage(w, r, id)
		case "PUT":
			putPage(w, r, id)
		case "DELETE":
			deletePage(w, r, id)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// Creates a new page with the specified title
func postPage(w http.ResponseWriter, r *http.Request) {
	var p Page
	if r.Body == nil {
		http.Error(w, "No Body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		logHttpError(w, r, err, "Error while parsing JSON")
		return
	}
	// Save page
	id, err := p.save()
	if err != nil {
		logHttpError(w, r, err, "Error while saving page")
		return
	}
	// Respond
	http.Redirect(w, r, string(id), http.StatusFound)
}
func getPage(w http.ResponseWriter, r *http.Request, id PageId) {
	// Load page
	p, err := loadPage(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Encode JSON
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		logHttpError(w, r, err, "Error while encoding JSON")
		return
	}
}
func putPage(w http.ResponseWriter, r *http.Request, id PageId) {
}
func deletePage(w http.ResponseWriter, r *http.Request, id PageId) {
}


func logHttpError(w http.ResponseWriter, r *http.Request, err error, message string) {
	log.Printf("%s: %s", message, err.Error())
	http.Error(w, message, http.StatusBadRequest)
}

type PageId string

type Page struct {
	Title string
}

func (p *Page) save() (PageId, error) {
	// Create id
	// Avoid characters that look the same: 0oO, lIi
	const chars = "bcdfghjkmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ" + "123456789"
	const strlen = 5
	randomBytes := make([]byte, strlen)
	for i := range randomBytes {
		randomBytes[i] = chars[rand.Intn(len(chars))]
	}
	id := PageId(randomBytes)

	// Save to map
	pages[id] = *p
	return id, nil
}

func loadPage(id PageId) (*Page, error) {
	p, present := pages[id]
	if !present {
		return nil, errors.New("id doesn't exist")
	}
	return &p, nil
}

var pages = make(map[PageId]Page);
