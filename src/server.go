package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	dbpassword string
	globalkey  = ""
	store      = sessions.NewCookieStore([]byte(globalkey))

	dbusername = "root"
	dbhostname = "mysql"
	dbname     = "test"
	// Global sql.DB to access the database by all handlers
	db         *sql.DB
	err        error
	seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

func home(w http.ResponseWriter, r *http.Request) {
	// confuse enumeration
	for i := 0; i < rand.Intn(50); i++ {
		fmt.Fprintln(w, fmt.Sprintf("<p hidden><a href='/%s'>%s</a></p>", randomString(rand.Intn(50)), randomString(rand.Intn(50))))
	}
	// show articles if get param set
	article, ok := r.URL.Query()["q"]

	if !ok || len(article[0]) < 1 {
		menu(w, r)
		serveFile(w, r, "html/homepage.html")

		// show articles

		files, err := ioutil.ReadDir("articles")
		if err != nil {
			fmt.Fprintln(w, err)
		}

		fmt.Fprintln(w, "<main><body><div>")
		for _, file := range files {
			fmt.Fprintln(w, "<a href='/home?q="+file.Name()+"'>"+file.Name()+"</a><br>")
		}
		fmt.Fprintln(w, "</main></body></div>")

		footer(w, r)
		return
	}

	// load article
	menu(w, r)

	// super safe sanitization techniques
	reg, err := regexp.Compile("../")
	if err != nil {
		log.Fatal(err)
	}
	filename := reg.ReplaceAllString(article[0], "")
	contents, err := ioutil.ReadFile(fmt.Sprintf("articles/%s", filename))
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintln(w, string(contents))

	footer(w, r)
}

func main() {

	// get .env secrets
	_ = godotenv.Load()
	dbpassword = os.Getenv("dbpassword")
	globalkey = os.Getenv("key")
	store = sessions.NewCookieStore([]byte(globalkey))

	// Create an sql.DB and check for errors
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true", dbusername, dbpassword, dbhostname, dbname))
	//db, err = sql.Open("mysql", fmt.Sprintf("%s@tcp(%s)/%s", dbusername, dbhostname, dbname))
	if err != nil {
		panic(err.Error())
	}
	// sql.DB should be long lived "defer" closes it once this function ends
	defer db.Close()

	// wait for the database to be ready when starting containers together
	ready := false
	for ready == false {
		// Test the connection to the database
		err = db.Ping()
		if err != nil {
			//panic(err.Error())
			time.Sleep(time.Second)
		} else {
			ready = true
		}
	}

	// set up handlers for endpoints
	http.HandleFunc("/home", logging(home, false))
	http.HandleFunc("/data", logging(data, true))
	http.HandleFunc("/send", logging(sendMessage, true))
	http.HandleFunc("/account", logging(account, true))
	http.HandleFunc("/signup", logging(signupPage, false))
	http.HandleFunc("/logout", logging(logoutPage, false))
	http.HandleFunc("/login", logging(login, false))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("html/css"))))
	http.HandleFunc("/", logging(home, true))

	// start the server
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// returns a random alphabetical string of provided length
func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
