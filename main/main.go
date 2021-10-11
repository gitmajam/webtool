package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

var tpl *template.Template
var dbSessions = map[string]string{}
var dbUsers = map[string]user{}
var db *sql.DB
var err error

const sessionLength int = 1000

// other form to declare an empty map
/*
var dbUsers = make(map[string]user)
*/

func init() {

	tpl = template.Must(template.ParseGlob("template/*"))

	//connection to data base azure
	//db, err = sql.Open("mysql", "awuser:Teflonio45@tcp(mydbinstance.cvbw65qn9mnn.us-east-2.rds.amazonaws.com)/test02?charset=utf8")

	//connection to data base local
	db, err = sql.Open("mysql", "root:Dexam3t@tcp(localhost:3306)/testgallery?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}

	//verify database connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("conexion a base de datos")
	}
}

func main() {

	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("public"))))
	http.Handle("/tmp/", http.StripPrefix("/tmp", http.FileServer(http.Dir("template"))))
	http.Handle("/", http.HandlerFunc(index))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/uploadGabriel", session(upGabriel))
	http.HandleFunc("/uploadDavid", session(upDavid))
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/gallery", session(gallery))
	http.HandleFunc("/controlPanel", session(permission(controlPanel)))
	http.HandleFunc("/logout", session(logout))

	http.HandleFunc("/controlPanel/users", users)
	// http.HandleFunc("/controlPanel/create", session(permission(create)))
	// http.HandleFunc("/controlPanel/insert", session(permission(insert)))
	// http.HandleFunc("/controlPanel/read", session(permission(read)))
	// http.HandleFunc("/controlPanel/update", session(permission(update)))
	// http.HandleFunc("/controlPanel/delete", session(permission(delete)))

	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {

	if sessionIsActive(res, req) {
		http.Redirect(res, req, "/gallery", http.StatusSeeOther)
		return
	}

	//verify login credentials
	if req.Method == http.MethodPost {

		username := req.FormValue("username")
		password := req.FormValue("password")

		//Compare the hashes
		err := bcrypt.CompareHashAndPassword(dbUsers[username].Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		if _, ok := dbUsers[username]; ok && password == password {

			//generate a cookie for the session
			cookie, err := req.Cookie("session")
			if err != nil {
				id := uuid.NewV4()
				cookie = &http.Cookie{
					Name:  "session",
					Value: id.String(),
					// Secure:   true,
					HttpOnly: true,
					// Path:     "/",
				}
				http.SetCookie(res, cookie)
			}

			//link the user with a session ID
			dbSessions[cookie.Value] = username

			http.Redirect(res, req, "/gallery", http.StatusSeeOther)
			return

		} else {

			//Executes template and goes out by res
			err := tpl.ExecuteTemplate(res, "index.gohtml", "Incorrect username or password!")
			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
			return
		}
	}

	//Executes template and goes out by res
	err := tpl.ExecuteTemplate(res, "index.gohtml", nil)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
}

func signup(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Entro a Signup")

	if req.Method == http.MethodPost {

		username := req.FormValue("username")
		password := req.FormValue("password")
		fname := req.FormValue("fname")
		lname := req.FormValue("lname")
		role := req.FormValue("role")

		// username taken?
		if _, ok := dbUsers[username]; ok {

			//Executes template and goes out by res
			err := tpl.ExecuteTemplate(res, "signup.gohtml", "Username already exist!")
			if err != nil {
				log.Fatalln("template didn't execute: ", err)
			}
			return

		} else {

			//encrypt password
			bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}
			//save user data
			dbUsers[username] = user{username, bs, fname, lname, role}
			fmt.Println("pass hash", bs)
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return
		}

	}

	//Executes template and goes out by res
	err := tpl.ExecuteTemplate(res, "signup.gohtml", nil)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}

	//This parse is only necesary if we're going to use req.Form
	/*
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
	*/
	// this call represents all the names and variables in a html Form map, req.ParseForm method should be applied before.

	/*
		for name, variable := range req.Form {
			fmt.Fprintln(res, name, variable)
		}
	*/

}

func gallery(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Entro a Gallery")

	cookie, err := req.Cookie("session")
	if err != nil {
		return
	}

	//update expire session time
	cookie.MaxAge = sessionLength
	http.SetCookie(res, cookie)

	if valuser, ok := getUser(req); ok {
		fmt.Println("Role:", valuser.Role)

		//Executes template and goes out by res
		err := tpl.ExecuteTemplate(res, "gallery.gohtml", valuser)
		if err != nil {
			log.Fatalln("template didn't execute: ", err)
		}
	} else {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
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
	http.Redirect(res, req, "/gallery", http.StatusSeeOther)
}

func logout(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Entro a Logout")

	cookie, err := req.Cookie("session")
	if err != nil {
		fmt.Println("There isn't cookie")
		return
	}

	cookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(res, cookie)
	http.Redirect(res, req, "/index", http.StatusSeeOther)
}

func controlPanel(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Entró a Contro Panel")

	u, _ := getUser(req)

	s := struct {
		Users map[string]user
		User  user
	}{
		Users: dbUsers,
		User:  u,
	}

	//Executes template and goes out by res
	err := tpl.ExecuteTemplate(res, "controlpanel.gohtml", s)
	if err != nil {
		log.Fatalln("template didn't execute: ", err)
	}
}

func upGabriel(res http.ResponseWriter, req *http.Request) {

	fmt.Println("Post Method received, form Gabriel")
	SaveFile(res, req, "./user/gabriel/")
}

func upDavid(res http.ResponseWriter, req *http.Request) {

	fmt.Println("Post Method received, form David")
	SaveFile(res, req, "./user/david/")

}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
