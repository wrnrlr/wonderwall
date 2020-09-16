package wonderwall

import (
	"fmt"
	"html/template"
	"net/http"
)

var (
	applicationJson = "application/json"
)

func ContentType(t string, w http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Content-Type") != t {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return false
	} else {
		return true
	}
}

func writeTmpl(w http.ResponseWriter, name string, i interface{}) {
	indexTmpl, err := template.ParseFiles(fmt.Sprintf("./template/%s.html", name))
	if err != nil {
		panic(err)
	}
	if err = indexTmpl.Execute(w, nil); err != nil {
		panic(err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	writeTmpl(w, "500", err)
	w.WriteHeader(500)
}

func RenderTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeTmpl(w, name, nil)
	}
}
