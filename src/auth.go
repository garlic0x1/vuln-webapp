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

func account(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if r.Method == "POST" {
		// modify status
		//username := r.FormValue("username") // can add idor here
		// super safe sanitization techniques
    		status := r.FormValue("status")
		
		result, err := db.Exec("UPDATE users SET status=? WHERE username=?", status, session.Values["username"])
		if err != nil {
			log.Println("ERROR UPDATING",err)
		}

		log.Println(result)

		http.Redirect(w, r, fmt.Sprintf("/account?user=%s", session.Values["username"]), 303)
	}
	menu(w, r)
	user, ok := r.URL.Query()["user"]

    	if !ok || len(user[0]) < 1 {
        	log.Println("Url Param 'key' is missing")
        	return
    	}

	// Grab from the database
	var databaseUsername string
	var databaseStatus string

	log.Println(user)
	// Search the database for the username provided
	// If it exists grab the password for validation
	err := db.QueryRow("SELECT username, status FROM users WHERE username=?", user[0]).Scan(&databaseUsername, &databaseStatus)
	// If not then redirect to the login page
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 303)
		return
	}

	// display username, status, etc
	fmt.Fprintln(w, "<body><main><div>")
	fmt.Fprintln(w, "<h2>Username</h2><big>")
	fmt.Fprintln(w, databaseUsername)
	fmt.Fprintln(w, "</big><h2>Status</h2><big>")
	fmt.Fprintln(w, databaseStatus)
	fmt.Fprintln(w, "</big>")
	fmt.Fprintln(w, "</body></main></div>")

	// if session user is the same, allow modification of status

	if session.Values["authenticated"] == true && session.Values["username"] == user[0] {
		serveFile(w, r, "html/update.html")
	}

	footer(w, r)
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
