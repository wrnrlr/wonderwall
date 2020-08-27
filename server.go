package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	indexHandle := func(w http.ResponseWriter, r *http.Request) {
		indexTmpl, err := template.ParseFiles("./template/index.html")
		if err != nil {
			panic(err)
		}
		err = indexTmpl.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", noCacheMiddleware(http.StripPrefix("/static/", fs).ServeHTTP))
	h := noCacheMiddleware(indexHandle)
	http.HandleFunc("/", h)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
