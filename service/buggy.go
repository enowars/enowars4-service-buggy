package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(w, "Hello world!")
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", mainHandler)
    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
