package routes

import (
	"buggy/go/db"
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type account struct {
	User     db.User
	Auth     bool
	Messages []db.Message
}

type productpage struct {
	Account  account
	Comments []db.Comment
}

type reg struct {
	Duplicate bool
	Error     bool
}

type login struct {
	Incorrect bool
	Error     bool
}

var store *sessions.CookieStore

var tpl *template.Template

func init() {
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(account{})

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// Index : Display main page
func Index(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	if acc.Auth {
		tpl.ExecuteTemplate(w, "index.gohtml", acc)
	} else {
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	}
}

// Register : Register new user
func Register(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	if acc.Auth {
		http.Redirect(w, req, "/", http.StatusFound)
	} else {
		if req.Method == http.MethodPost {
			username := req.FormValue("username")
			password := req.FormValue("pw")
			if username != "" && password != "" {
				insert := db.InsertUser(username, password, "", false)
				if insert {
					sendWelcome(username)
					redirectOnSuccess(username, session, w, req)
				} else {
					tpl.ExecuteTemplate(w, "register.gohtml", reg{true, false})
				}
			} else {
				tpl.ExecuteTemplate(w, "register.gohtml", reg{false, true})
			}
		} else {
			tpl.ExecuteTemplate(w, "register.gohtml", nil)
		}
	}
}

// Login user
func Login(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	if acc.Auth {
		http.Redirect(w, req, "/", http.StatusFound)
	} else {
		if req.Method == http.MethodPost {
			username := req.FormValue("username")
			password := req.FormValue("pw")
			if username != "" && password != "" {

				loginValid := db.AuthUser(username, password)
				if loginValid {
					redirectOnSuccess(username, session, w, req)
				} else {
					err = session.Save(req, w)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					tpl.ExecuteTemplate(w, "login.gohtml", login{true, false})
				}
			} else {
				tpl.ExecuteTemplate(w, "login.gohtml", login{false, true})
			}
		} else {
			tpl.ExecuteTemplate(w, "login.gohtml", nil)
		}
	}

}

// Logout user
func Logout(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["account"] = account{}
	session.Options.MaxAge = -1

	err = session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

// Profile : Show user profile
func Profile(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Hot reload
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	acc := getAccount(session)
	messages := db.GetMessages(acc.User.Username)
	acc.Messages = messages
	if acc.Auth {
		tpl.ExecuteTemplate(w, "profile.gohtml", acc)
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func getAccount(s *sessions.Session) account {
	val := s.Values["account"]
	var acc = account{}
	acc, ok := val.(account)
	if !ok {
		return account{Auth: false}
	}
	return acc
}

func redirectOnSuccess(username string, session *sessions.Session, w http.ResponseWriter, req *http.Request) {
	user := db.GetUser(username)
	accountAuth := &account{
		User: user,
		Auth: true,
	}

	session.Values["account"] = accountAuth

	err := session.Save(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

func sendWelcome(username string) {
	db.AddMessage(username, "buggy-team", "Welcome to the one and only Buggy Store, enjoy your stay!")
}

func sendPreorder(username string, buggy string) {
	db.AddMessage(username, "buggy-team", fmt.Sprintf("Thank you for preordering the %s! We will inform you when it becomes available ASAP.", buggy))
}

// ProductOne : Product page for super buggy
func ProductOne(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	page := productpage{}
	page.Account = acc
	comments := db.GetComments("super-buggy")
	page.Comments = comments
	if acc.Auth {
		if req.Method == http.MethodPost {
			req.ParseForm()
			if req.Form["comment"] != nil {
				db.AddComment(acc.User.Username, "super-buggy", req.FormValue("comment"))
				comments := db.GetComments("super-buggy")
				page.Comments = comments
			} else {
				sendPreorder(acc.User.Username, "Super Buggy")
			}
			http.Redirect(w, req, "/super-buggy", http.StatusFound)
		} else {
			tpl.ExecuteTemplate(w, "super-buggy.gohtml", page)
		}
	} else {
		tpl.ExecuteTemplate(w, "super-buggy.gohtml", page)
	}
}

// ProductTwo : Product page for mega buggy
func ProductTwo(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	page := productpage{}
	page.Account = acc
	comments := db.GetComments("mega-buggy")
	page.Comments = comments
	if acc.Auth {
		if req.Method == http.MethodPost {
			req.ParseForm()
			if req.Form["comment"] != nil {
				db.AddComment(acc.User.Username, "mega-buggy", req.FormValue("comment"))
				comments := db.GetComments("mega-buggy")
				page.Comments = comments
			} else {
				sendPreorder(acc.User.Username, "Mega Buggy")
			}
			http.Redirect(w, req, "/mega-buggy", http.StatusFound)
		} else {
			tpl.ExecuteTemplate(w, "mega-buggy.gohtml", page)
		}
	} else {
		tpl.ExecuteTemplate(w, "mega-buggy.gohtml", page)
	}
}
