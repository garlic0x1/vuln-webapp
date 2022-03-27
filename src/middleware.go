package main

import (
    "log"
    "net/http"
)

func logging(f http.HandlerFunc, needauth bool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

	    // if authorization is required, check that the client's cookies and if
	    // they aren't logged in, redirect them to the login page
	    if needauth {
		session, _ := store.Get(r, "session")

    		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        		http.Redirect(w, r, "/login", http.StatusSeeOther)
        		return
    		}
	    }
	    
	    // the rest of the logging is done in caddy
      	    log.Println(r.URL, r.Header.Get("X-Forwarded-For"))
            f(w, r)
    }
}

