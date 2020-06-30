package main

import (
	"buggy/go/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates/"))))

	router.HandleFunc("/", routes.Index)
	router.HandleFunc("/register", routes.Register)
	router.HandleFunc("/login", routes.Login)
	router.HandleFunc("/logout", routes.Logout)
	router.HandleFunc("/profile", routes.Profile)
	router.HandleFunc("/tickets", routes.Ticket)
	router.HandleFunc("/tickets/{hash}", routes.Tickets)

	router.HandleFunc("/super-buggy", routes.ProductOne)
	router.HandleFunc("/mega-buggy", routes.ProductTwo)

	log.Fatal(http.ListenAndServe(":7890", router))
}
