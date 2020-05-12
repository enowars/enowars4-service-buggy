package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type User struct {
	Username      string
	Authenticated bool
	Admin         bool
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
	// TODO: Implement these
	// router.HandleFunc("/login", login)
	// router.HandleFunc("/logout", logout)
	// router.HandleFunc("/register", logout)
	// router.HandleFunc("/forbidden", forbidden)
	// router.HandleFunc("/profile", profile)

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// First vuln (?), predictable admin cookie
func initCookies() {
	//rand.Seed(<team_id>)
	var adminCookie = rand.Intn(1000000)
	fmt.Printf("Using admin cookie %d\n", adminCookie)
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}
