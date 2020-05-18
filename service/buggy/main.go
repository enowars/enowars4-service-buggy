package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type user struct {
	Username string
	Password string
	Status   string
	Admin    bool
}

var store *sessions.CookieStore

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	router := mux.NewRouter()
	// Serve images
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	// Servce html/css
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))

	router.HandleFunc("/", index)
	router.HandleFunc("/profile", profile)
	// TODO: Implement these
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/register", register)
	// router.HandleFunc("/forbidden", forbidden)

	fmt.Println("Running server on port 7890")
	log.Fatal(http.ListenAndServe(":7890", router))
}

// First vuln (?), predictable admin cookie
func initCookies() {
	//rand.Seed(<team_id>)
	var adminCookie = rand.Intn(1000000)
	fmt.Printf("Using admin cookie %d\n", adminCookie)
}

// TODO: Templating doesn't fully work yet
func register(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if req.FormValue("username") != "" && req.FormValue("pw") != "" {
			insert := insertUser(req.FormValue("username"), req.FormValue("password"), "", false)
			if insert {
				tpl.ExecuteTemplate(w, "register.gohtml", struct{ Success bool }{Success: true})
			} else {
				tpl.ExecuteTemplate(w, "register.gohtml", struct{ Duplicate bool }{Duplicate: true})
			}
		} else {
			tpl.ExecuteTemplate(w, "register.gohtml", struct{ Error bool }{Error: true})
		}
	} else {
		tpl.ExecuteTemplate(w, "register.gohtml", nil)
	}
}

// TODO: Implement all routes

func login(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func logout(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func profile(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func insertUser(username string, pw string, status string, admin bool) bool {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/enodb")

	if err != nil {
		return false
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE name = ?)", username).Scan(&exists)

	if err != nil || exists {
		return false
	}

	insert, err := db.Query("INSERT IGNORE INTO users VALUES (?, ?, ?, ?)", username, pw, status, admin)
	if err != nil {
		return false
	}
	defer insert.Close()

	return true
}

func authUser(username string, pw string) bool {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/enodb")

	if err != nil {
		return false
	}
	defer db.Close()

	var userReq user
	err = db.QueryRow("SELECT name, password FROM users WHERE BINARY name = ?", username).Scan(&userReq.Username, &userReq.Password)

	// Ohne "BINARY" case-insensitiv -> vuln falls man zb Account mit username==admIN erstellen kann
	// Wollen wir das als leichte vuln?
	err = db.QueryRow("SELECT name, password FROM users WHERE name = ?", username).Scan(&userReq.Username, &userReq.Password)
	if err != nil {
		return false
	}

	if pw != userReq.Password {
		return false
	}

	return true
}

func deleteUser(username string) bool {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/enodb")

	if err != nil {
		return false
	}
	defer db.Close()

	delete, err := db.Query("DELETE FROM users WHERE BINARY name = ?", username)
	if err != nil {
		return false
	}
	defer delete.Close()

	return true
}

// TODO: Remove at some point, only to be used now for debugging
func printDB() {

	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/enodb")

	if err != nil {
		log.Fatal("Error connecting to database.")
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM users")

	if err != nil {
		log.Fatal("Query error.")
	}

	for results.Next() {
		var userTest user

		err = results.Scan(&userTest.Username, &userTest.Password, &userTest.Status, &userTest.Admin)

		if err != nil {
			log.Printf("noo")
		}

		log.Printf("%s %s %s %t\n", userTest.Username, userTest.Password, userTest.Status, userTest.Admin)
	}
}
