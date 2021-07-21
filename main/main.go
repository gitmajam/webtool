package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	UserName string
	First    string
	Last     string
}

var tpl *template.Template
var dbSessions = map[string]string{}
var dbUsers = map[string]user{}

// other form to declare an empty map
/*
var dbUsers = make(map[string]user)
*/

func init() {
	tpl = template.Must(template.ParseFiles("template/index.gohtml"))
}

func main() {

	dbUsers["marlonjavi@gmail.com"] = user{"marlonjavi@gmail.com", "muñoz", "uribe"}
	fmt.Println(dbUsers["marlonjavi@gmail.com"])

	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("public"))))
	http.Handle("/tmp/", http.StripPrefix("/tmp", http.FileServer(http.Dir("template"))))
	http.Handle("/", http.HandlerFunc(a))
	http.Handle("/uploadGabriel", http.HandlerFunc(upGabriel))
	http.Handle("/uploadDavid", http.HandlerFunc(upDavid))
	http.Handle("/read", http.HandlerFunc(read))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/login", http.HandlerFunc(login))

	http.ListenAndServe("Localhost:8080", nil)
}

func a(res http.ResponseWriter, req *http.Request) {

	//generate a cookie for each session
	cookie, err := req.Cookie("session-id")
	if err != nil {
		id := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session-id",
			Value: id.String(),
			// Secure:   true,
			HttpOnly: true,
			// Path:     "/",
		}
		http.SetCookie(res, cookie)
	}

	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	fmt.Println(req.URL)
	fmt.Println(req.Method)

	//Executes the template and goes out by res
	err = tpl.Execute(res, nil)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
}

func login(res http.ResponseWriter, req *http.Request) {

	fmt.Println(dbUsers["marlonjavi@gmail.com"])

	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	//read cookie
	cookie, err := req.Cookie("session-id")
	if err != nil {
		fmt.Fprintln(res, "No Cookie provided")
	}

	//if the user exists then link the username with sessionID, if not, register the user and link it

	username := req.FormValue("username")
	fname := req.FormValue("fname")
	lname := req.FormValue("lname")

	if _, found := dbSessions[cookie.Value]; !found {
		if _, found := dbUsers[username]; !found {
			dbUsers[username] = user{username, fname, lname}
		}
		dbSessions[cookie.Value] = username
	}

	u := dbSessions[cookie.Value]
	data := dbUsers[u]

	err = tpl.Execute(res, data)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}

	//fmt.Fprintln(res, dbUsers[dbSessions[cookie.Value]])

	// this call represents all the names and variables in a html Form map, req.ParseForm method should be applied before.
	/*
		for name, variable := range req.Form {
			fmt.Fprintln(res, name, variable)
		}
	*/

}

//shows the cookie name and ID in the browser
func read(res http.ResponseWriter, req *http.Request) {

	c, err := req.Cookie("session-id")
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
		return
	}

	err = tpl.ExecuteTemplate(res, "index.gohtml", nil)
	HandleError(res, err)
}

func HandleError(res http.ResponseWriter, err error) {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}
