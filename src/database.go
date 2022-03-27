package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

func searchPage(res http.ResponseWriter, req *http.Request) {
	search, err := template.ParseFiles("html/search.html")
	if err != nil {
		fmt.Fprintln(res, err)
	}
	err = search.Execute(res, nil)
	if err != nil { // if there is an error
		fmt.Fprintln(res, err)
		log.Println("template executing error: ", err) //log it
	}
	footer(res, req)
}

func data(res http.ResponseWriter, req *http.Request) {
	menu(res, req)

	// If method is GET serve an query page
	if req.Method != "POST" {
		searchPage(res, req)
		return
	}

	// safety first
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		fmt.Fprintln(res, err)
		log.Fatal(err)
	}
	_ = reg

	// get post data to build query
	quser := fmt.Sprintf(`'%%%s%%'`, req.FormValue("user"))
	// oops forgot the regex

	buildquery := fmt.Sprintf("select id, username, status from users where username like %s", quser)
	log.Println(buildquery)

	// Grab from the database
	rows, err := db.Query(buildquery)
	if err != nil {
		fmt.Fprintln(res, err)
		fmt.Println("error performing query", err)
		return
	}
	defer rows.Close()

	//html table
	type tpl struct {
		ID     int
		User   string
		Status string
	}

	type Tpl struct {
		Users []tpl
	}

	var build_tpls []tpl

	// populate the template struct
	for rows.Next() {
		var id int
		var user string
		var status string

		err = rows.Scan(&id, &user, &status)
		if err != nil {
			// handle this error
			panic(err)
		}
		build_tpls = append(build_tpls, tpl{
			ID:     id,
			User:   user,
			Status: status,
		})
	}
	form, err := template.ParseFiles("html/users.html")
	if err != nil {
		fmt.Fprintln(res, err)
	}
	err = form.Execute(res, Tpl{
		Users: build_tpls,
	})
	if err != nil {
		log.Println("Error executing template", err)
	}
		footer(res, req)
}
