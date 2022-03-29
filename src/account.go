package main

import (
	"fmt"
	"html/template"
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
	var userID int

	// Search the database for the username provided
	err := db.QueryRow("SELECT id, username, status FROM users WHERE username=?", user[0]).Scan(&userID, &databaseUsername, &databaseStatus)
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
	// and show messages
	if session.Values["authenticated"] == true && session.Values["username"] == user[0] {

		serveFile(w, r, "html/update.html")
		var dbID int
		err = db.QueryRow("select id from users where username=?", user[0]).Scan(&dbID)
		if err != nil {
			fmt.Fprintln(w, "Error performing query", err)
			return
		}

		rows, err := db.Query("select sender, reciever, message from messages where reciever=?", dbID)
		if err != nil {
			fmt.Fprintln(w, "Error performing query", err)
			return
		}
		rows2, err := db.Query("select sender, reciever, message from messages where sender=?", dbID)
		if err != nil {
			fmt.Fprintln(w, "Error performing query", err)
			return
		}

		defer rows2.Close()
		defer rows.Close()

		//html table
		type tpl struct {
			Sender     string
			Reciever   string
			Message string
		}
	
		type Tpl struct {
			Sent []tpl
			Recieved []tpl
		}
	
		var build_recieved []tpl
		var build_sent []tpl
	
		// populate the template struct
		for rows2.Next() {
			var senderID int
			var senderName string
			var recieverID int
			var recieverName string
			var message string
	
			err = rows2.Scan(&senderID, &recieverID, &message)
			if err != nil {
				// handle this error
				panic(err)
			}

			// get usernames from ID
			err = db.QueryRow("select username from users where id=?", senderID).Scan(&senderName)
			if err != nil { log.Println(err) }
			err = db.QueryRow("select username from users where id=?", recieverID).Scan(&recieverName)
			if err != nil { log.Println(err) }
			build_sent = append(build_sent, tpl{
				Sender:     senderName,
				Reciever:   recieverName,
				Message:    message,
			})
		}
		for rows.Next() {
			var senderID int
			var senderName string
			var recieverID int
			var recieverName string
			var message string
	
			err = rows.Scan(&senderID, &recieverID, &message)
			if err != nil {
				// handle this error
				panic(err)
			}

			// get usernames from ID
			err = db.QueryRow("select username from users where id=?", senderID).Scan(&senderName)
			err = db.QueryRow("select username from users where id=?", recieverID).Scan(&recieverName)

			build_recieved = append(build_recieved, tpl{
				Sender:     senderName,
				Reciever:   recieverName,
				Message:    message,
			})
		}
		form, err := template.ParseFiles("html/messages.html")
		if err != nil {
			fmt.Fprintln(w, err)
		}
		err = form.Execute(w, Tpl{
			Sent: build_sent,
			Recieved: build_recieved,
		})
		if err != nil {
			log.Println("Error executing template", err)
		}
	} else {
		// if its another user allow sending a message
		var senderID int
		var recieverID int
		senderName := session.Values["username"]
		recieverName := user[0]

		// get visitors ID
		

		err = db.QueryRow("select id from users where username=?", senderName).Scan(&senderID)
		err = db.QueryRow("select id from users where username=?", recieverName).Scan(&recieverID)


		type tpl struct {
			SenderID     int
			RecieverID   int
		}
		form, err := template.ParseFiles("html/messageform.html")
		if err != nil {
			fmt.Fprintln(w, err)
		}
		err = form.Execute(w, tpl{
				SenderID: senderID,
				RecieverID: recieverID,
			})
		if err != nil {
			log.Println("Error executing template", err)
		}
	}

	footer(w, r)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if r.Method == "POST" {
		// send message
		//username := r.FormValue("username") // can add idor here
    		message := r.FormValue("message")
		senderID := r.FormValue("sender")
		recieverID := r.FormValue("reciever")
		
		_, err = db.Exec("INSERT INTO messages(sender, reciever, message) VALUES(?, ?, ?)", senderID, recieverID, message)
		if err != nil {
			fmt.Println(w, "Server error, unable to create your account.", 500)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/account?user=%s", session.Values["username"]), 303)
	}
}
