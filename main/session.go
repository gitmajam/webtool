package main

import (
	"fmt"
	"net/http"
)

func getUser(req *http.Request) (user, bool) {
	var u user

	//get the cookie session
	cookie, err := req.Cookie("session")
	if err != nil {
		fmt.Println("There isn't cookie")
		return u, false
	}

	// verify if the user has a session registered
	if un, ok := dbSessions[cookie.Value]; ok {
		u, ok = dbUsers[un]
		return u, ok
	} else {
		return u, false
	}
}

func sessionIsActive(res http.ResponseWriter, req *http.Request) bool {

	//check if there is a cookie
	cookie, err := req.Cookie("session")
	if err != nil {
		return false
	}

	//update expire session time
	cookie.MaxAge = sessionLength
	http.SetCookie(res, cookie)

	fmt.Println("ID sesion:", cookie.Value)

	//check if user is already register
	if un, ok := dbSessions[cookie.Value]; ok {
		_, ok = dbUsers[un]
		return ok
	} else {
		return false
	}
}
