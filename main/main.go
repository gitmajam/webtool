package main

import (
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("template/index.html"))
}

func a(res http.ResponseWriter, req *http.Request) {
	err := tpl.Execute(res, nil)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}

}

func main() {

	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("public"))))
	http.Handle("/tmp/", http.StripPrefix("/tmp", http.FileServer(http.Dir("template"))))
	http.Handle("/", http.HandlerFunc(a))
	http.ListenAndServe("Localhost:8080", nil)

}
