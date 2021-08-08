package main

import (
	"fmt"
	"net/http"
)

func users(res http.ResponseWriter, req *http.Request) {

	rows, err := db.Query(`SELECT userName FROM users;`)
	check(err)
	defer rows.Close()
	fmt.Println("hace el query")

	//data to be use in query
	var s, name string
	s = "RETRIEVE RECORD:\n"

	//query
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		s += name + "\n"
	}
	fmt.Fprintln(res, s)
}

func create(res http.ResponseWriter, req *http.Request) {

}

func insert(res http.ResponseWriter, req *http.Request) {

}

func read(res http.ResponseWriter, req *http.Request) {

}

func update(res http.ResponseWriter, req *http.Request) {

}

func delete(res http.ResponseWriter, req *http.Request) {

}
