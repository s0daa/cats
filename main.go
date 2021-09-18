package main

import (
	"errors"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Mimetypes = map[string]bool{".png": true, ".jpg": true, ".jpeg": true}

var images []string

func main() {
	const p = "."
	images = findInPath(p)
	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/cats/"+getRandomCat(), 321)
	}
	http.HandleFunc("/", indexHandler)

	fileServer := http.FileServer(imgDir(p))
	http.Handle("/cats/", http.StripPrefix("/cats/", fileServer))

	catHandler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		http.ServeFile(w, req, getRandomCat())
	}
	http.HandleFunc("/cat", catHandler)

	iconHandler := func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "./favicon.svg")
	}
	http.HandleFunc("/favicon.ico", iconHandler)

	log.Fatal(http.ListenAndServe(":8090", nil))
}

type imgDir string

func (d imgDir) Open(name string) (http.File, error) {
	if !Mimetypes[strings.ToLower(filepath.Ext(name))] {
		return nil, errors.New("not image")
	}
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("invalid character in file path")
	}
	dir := string(d)
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	f, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func findInPath(path string) []string {
	var a []string
	filepath.WalkDir(path, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if Mimetypes[filepath.Ext(d.Name())] {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func getRandomCat() string {
	return images[rand.Intn(len(images)-0)+0]
}
