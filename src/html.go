package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type menuVars struct {
	Account  string
	User     string
	Inouturl string
	Inout    string
}

func menu(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "session")
	if err != nil {
		fmt.Println("error getting session")
	}

	// create data variable to be given to the html template
	// depending on the users authentication status they will
	// see a different menu
	var d menuVars
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// present login/sign up
		d = menuVars{
			Account:  "/signup",
			User:     "Sign Up",
			Inouturl: "/login",
			Inout:    "Login",
		}
	} else {
		// present logout/name
		if user, ok := session.Values["username"].(string); ok {
			d = menuVars{
				Account:  fmt.Sprintf("/account?user=%s", user),
				User:     user,
				Inouturl: "/logout",
				Inout:    "Logout",
			}
		}
	}

	// parse and and execute the template
	tmpl, err := template.ParseFiles("html/menu.html")
	if err != nil {
		fmt.Fprintln(w, err)
	}
	err = tmpl.Execute(w, d)
	if err != nil {
		fmt.Println(err)
	}
}

func serveFile(w http.ResponseWriter, req *http.Request, file string) {
	footer, err := template.ParseFiles(file)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	err = footer.Execute(w, nil)
	if err != nil { // if there is an error
		log.Println("template executing error: ", err) //log it
	}
}

func footer(w http.ResponseWriter, req *http.Request) {
	serveFile(w, req, "html/footer.html")
}
