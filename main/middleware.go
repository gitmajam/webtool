package main

import (
	"net/http"
)

func session(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !sessionIsActive(res, req) {
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(res, req)
	})
}

func permission(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//redirect to "/" if there isn't user or user isn't admin role
		if valuser, ok := getUser(req); !ok || valuser.Role != "superadmin" {
			http.Redirect(res, req, "/", http.StatusSeeOther)
		}
		h.ServeHTTP(res, req)
	})
}
