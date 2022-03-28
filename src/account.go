package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func account(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if r.Method == "POST" {
		// modify status
		//username := r.FormValue("username") // can add idor here
    		status := r.FormValue("status")
		
		_, err := db.Exec("UPDATE users SET status=? WHERE username=?", status, session.Values["username"])
		if err != nil {
			log.Println("ERROR UPDATING",err)
		}

		http.Redirect(w, r, fmt.Sprintf("/account?user=%s", session.Values["username"]), 303)
	}
	menu(w, r)
	user, ok := r.URL.Query()["user"]

    	if !ok || len(user[0]) < 1 {
        	log.Println("Url Param 'user' is missing")
        	return
    	}

	// Grab from the database
	var databaseUsername string
	var databaseStatus string

	// Search the database for the username provided
	err := db.QueryRow("SELECT username, status FROM users WHERE username=?", user[0]).Scan(&databaseUsername, &databaseStatus)
	// If not exists then redirect to the home page
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", 303)
		return
	}

	// display username, status, etc
	fmt.Fprintln(w, "<body><main><div>")
	fmt.Fprintln(w, "<a href='javascript:window.history.back()'>Back</a><br><br>")
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
