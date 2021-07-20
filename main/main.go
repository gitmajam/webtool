package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("template/index.gohtml"))
}

func main() {

	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("public"))))
	http.Handle("/tmp/", http.StripPrefix("/tmp", http.FileServer(http.Dir("template"))))
	http.Handle("/", http.HandlerFunc(a))
	http.Handle("/uploadGabriel", http.HandlerFunc(upGabriel))
	http.Handle("/uploadDavid", http.HandlerFunc(upDavid))
	http.Handle("/read", http.HandlerFunc(read))
	http.ListenAndServe("Localhost:8080", nil)

}

func a(res http.ResponseWriter, req *http.Request) {

	//cookie
	http.SetCookie(res, &http.Cookie{
		Name:  "my-cookie",
		Value: "14041980",
		Path:  "/",
	})
	fmt.Println("Cookie written, check your browser")

	//Executes the template and goes out by res
	err := tpl.Execute(res, nil)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
}
func read(res http.ResponseWriter, req *http.Request) {

	c, err := req.Cookie("my-cookie")
	if err != nil {
		http.Error(res, http.StatusText(400), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(res, "Your cookie:", c)
}

func upGabriel(res http.ResponseWriter, req *http.Request) {

	fmt.Println("Post Method received, form Gabriel")
	SaveFile(res, req, "./user/gabriel/")
}

func upDavid(res http.ResponseWriter, req *http.Request) {

	fmt.Println("Post Method received, form David")
	SaveFile(res, req, "./user/david/")

}

func SaveFile(res http.ResponseWriter, req *http.Request, directory string) {

	//Open the File received
	f, h, err := req.FormFile("q")

	//error handling
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close() //close the File when it's not longer used

	//read
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	//store on server
	dst, err := os.Create(filepath.Join(directory, h.Filename))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = dst.Write(bs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	error := tpl.ExecuteTemplate(res, "index.gohtml", nil)
	HandleError(res, error)
}

func HandleError(res http.ResponseWriter, err error) {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
