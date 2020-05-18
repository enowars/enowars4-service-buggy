package routes

import (
	"buggy/go/db"
	"html/template"
	"net/http"
)

type reg struct {
	Success   bool
	Duplicate bool
	Error     bool
}

var tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

// Index : Display main page
func Index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

// Register : Register new user
func Register(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if req.FormValue("username") != "" && req.FormValue("pw") != "" {
			insert := db.InsertUser(req.FormValue("username"), req.FormValue("password"), "", false)
			if insert {
				tpl.ExecuteTemplate(w, "register.gohtml", reg{true, false, false})
			} else {
				tpl.ExecuteTemplate(w, "register.gohtml", reg{false, true, false})
			}
		} else {
			tpl.ExecuteTemplate(w, "register.gohtml", reg{false, false, true})
		}
	} else {
		tpl.ExecuteTemplate(w, "register.gohtml", nil)
	}
}

// TODO: Implement all routes

// Login user
func Login(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

// Logout user
func Logout(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

// Profile : Show user profile
func Profile(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}
