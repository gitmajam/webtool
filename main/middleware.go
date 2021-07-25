package main

import (
	"net/http"
)

func authorized(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !sessionIsActive(res, req) {
			http.Redirect(res, req, "/", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(res, req)
	})
}
