package routes

import (
	"buggy/go/db"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type account struct {
	User     db.User
	Auth     bool
	Messages []db.Message
	Tickets  []db.Ticket
}

type productpage struct {
	Account  account
	Comments []db.Comment
}
type ticketpage struct {
	Account account
	Ticket  db.Ticket
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

func ProductOne(w http.ResponseWriter, req *http.Request) {
	productPage(w, req, "super")
}

func ProductTwo(w http.ResponseWriter, req *http.Request) {
	productPage(w, req, "mega")
}

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

func Profile(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	messages := db.GetMessages(acc.User.Username)
	tickets := db.GetTickets(acc.User.Username)
	acc.Messages = messages
	acc.Tickets = tickets
	if acc.Auth {
		tpl.ExecuteTemplate(w, "profile.gohtml", acc)
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func Ticket(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	if acc.Auth {
		if req.Method == http.MethodPost {
			subject := req.FormValue("subject")
			message := req.FormValue("message")
			if subject != "" && message != "" {
				str := acc.User.Username + strconv.FormatInt(time.Now().Unix(), 10)
				sha := sha256.Sum256([]byte(str))
				hash := hex.EncodeToString(sha[:])
				db.AddMessage("buggy-team", acc.User.Username, hash, message)
				db.AddTicket(acc.User.Username, subject, hash)
				db.AddMessage(acc.User.Username, "buggy-team", hash, "Please be aware that the Buggy Store(tm) Team is really busy right now. Replies might be delayed.")
				http.Redirect(w, req, fmt.Sprintf("/tickets/%s", hash), http.StatusFound)
			} else {
				tpl.ExecuteTemplate(w, "ticket.gohtml", acc)
			}
		} else {
			tpl.ExecuteTemplate(w, "ticket.gohtml", acc)
		}
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func Tickets(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(req)
	hash := vars["hash"]
	account := getAccount(session)
	if account.Auth {
		if len(hash) == 64 {
			messages := db.GetAllMessages(hash)
			if len(messages) < 1 {
				http.Redirect(w, req, "/", http.StatusFound)
			} else {
				account.Messages = messages
				page := ticketpage{}
				page.Account = account
				page.Ticket = db.GetTicket(hash)
				tpl.ExecuteTemplate(w, "tickets.gohtml", page)
			}
		} else {
			http.Redirect(w, req, "/", http.StatusFound)
		}
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
	db.AddMessage(username, "buggy-team", "private", "Welcome to the one and only Buggy Store, enjoy your stay!")
}

func sendPreorder(username string, buggy string) {
	db.AddMessage(username, "buggy-team", "private", fmt.Sprintf("Thank you for preordering the %s! We will inform you when it becomes available ASAP.", buggy))
}

func productPage(w http.ResponseWriter, req *http.Request, buggy string) {
	buggy = strings.ToLower(buggy)

	session, err := store.Get(req, "buggy-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	acc := getAccount(session)
	page := productpage{}
	page.Account = acc
	comments := db.GetComments(buggy + "-buggy")
	page.Comments = sortComments(comments, acc.User.Username)
	if acc.Auth {
		if req.Method == http.MethodPost {
			req.ParseForm()
			if req.Form["comment"] != nil {
				db.AddComment(acc.User.Username, buggy+"-buggy", req.FormValue("comment"))
				comments := db.GetComments(buggy + "-buggy")
				page.Comments = comments
			} else {
				sendPreorder(acc.User.Username, strings.Title(buggy)+" Buggy")
			}
			http.Redirect(w, req, "/"+buggy+"-buggy", http.StatusFound)
		} else {
			tpl.ExecuteTemplate(w, buggy+"-buggy.gohtml", page)
		}
	} else {
		tpl.ExecuteTemplate(w, buggy+"-buggy.gohtml", page)
	}
}

func sortComments(comments []db.Comment, username string) []db.Comment {
	commentsOther := []db.Comment{}
	commentsUser := []db.Comment{}
	for _, c := range comments {
		if c.User == username {
			commentsUser = append(commentsUser, c)
		} else {
			commentsOther = append(commentsOther, c)
		}
	}
	return append(commentsUser, commentsOther...)
}
