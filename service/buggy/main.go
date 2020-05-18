package main

import (
	"buggy/go/routes"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	// Serve images
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	// Servce html/css
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))

	router.HandleFunc("/", routes.Index)
	router.HandleFunc("/register", routes.Register)
	router.HandleFunc("/login", routes.Login)
	router.HandleFunc("/logout", routes.Logout)
	router.HandleFunc("/profile", routes.Profile)

	fmt.Println("Running server on port 7890")
	log.Fatal(http.ListenAndServe(":7890", router))
}

// First vuln (?), predictable admin cookie
func initCookies() {
	//rand.Seed(<team_id>)
	var adminCookie = rand.Intn(1000000)
	fmt.Printf("Using admin cookie %d\n", adminCookie)
}
