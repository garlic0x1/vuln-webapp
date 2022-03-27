package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func login(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	// If method is GET serve an html login page
	if req.Method != "POST" {
		menu(res, req)

		serveFile(res, req, "html/login.html")

		footer(res, req)
		return
	}
	// Grab the username/password from the submitted post form
	username := req.FormValue("username")
	password := req.FormValue("password")

	// Grab from the database
	var databaseUsername string
	var databasePassword string

	// Search the database for the username provided
	// If it exists grab the password for validation
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
	// If not then redirect to the login page
	if err != nil {
		http.Redirect(res, req, "/login", 303)
		return
	}

	// Validate the password
	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	// If wrong password redirect to the login
	if err != nil {
		http.Redirect(res, req, "/login", 303)
		return
	}

	// If the login succeeded
	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Values["username"] = databaseUsername
	session.Save(req, res)
	http.Redirect(res, req, "/", 303)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")

	// Serve signup.html to get requests to /signup
	if req.Method != "POST" {

		menu(res, req)

		serveFile(res, req, "html/signup.html")

		footer(res, req)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	username = reg.ReplaceAllString(username, "")
	password = reg.ReplaceAllString(password, "")

	var usertemp string

	err = db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&usertemp)

	if len(username) >= 6 && len(password) >= 6 {
		switch {
		// Username is available
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Println(res, "Server error, unable to create your account.", 500)
				return
			}

			_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
			if err != nil {
				fmt.Println(res, "Server error, unable to create your account.", 500)
				return
			}
			//res.Write([]byte("User created!"))
			session.Values["authenticated"] = true
			session.Values["username"] = username
			session.Save(req, res)
			http.Redirect(res, req, "/", 303)
			return
		case err != nil:
			http.Error(res, "Server error, unable to create your account.", 500)
			fmt.Println("Error performing query", err)
			return
		default:
			http.Redirect(res, req, "/", 303)
		}
	} else {
		http.Redirect(res, req, "/signup", 303)
	}
}
