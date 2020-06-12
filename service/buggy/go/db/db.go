package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// User struct
type User struct {
	Username string
	Password string
	Status   string
	Admin    bool
}

// Message struct
type Message struct {
	To      string
	From    string
	Hash    string
	Content string
}

// Comment struct
type Comment struct {
	User    string
	Product string
	Content string
}

// Ticket struct
type Ticket struct {
	User    string
	Subject string
	Hash    string
}

// InsertUser : Insert user if not present
func InsertUser(username string, pw string, status string, admin bool) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

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

// AuthUser : Authenticate user
func AuthUser(username string, pw string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return false
	}
	defer db.Close()

	var userReq User
	err = db.QueryRow("SELECT name, password FROM users WHERE BINARY name = ?", username).Scan(&userReq.Username, &userReq.Password)

	// Ohne "BINARY" case-insensitiv -> vuln falls man zb Account mit username==admIN erstellen kann
	// Wollen wir das als leichte vuln?
	// err = db.QueryRow("SELECT name, password FROM users WHERE name = ?", username).Scan(&userReq.Username, &userReq.Password)
	if err != nil {
		return false
	}

	if pw != userReq.Password {
		return false
	}

	return true
}

// DeleteUser : Delete user if present
func DeleteUser(username string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

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

// GetUser : Return User from db if existing
func GetUser(username string) User {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return User{}
	}
	defer db.Close()

	var userReq User
	err = db.QueryRow("SELECT name, password, status, admin FROM users WHERE name = ?", username).Scan(&userReq.Username, &userReq.Password, &userReq.Status, &userReq.Admin)

	if err != nil {
		return User{}
	}
	return userReq
}

// AddMessage : Add message
func AddMessage(username string, sender string, hash string, content string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return false
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO messages VALUES (?, ?, ?, ?)", username, sender, hash, content)
	if err != nil {
		return false
	}
	defer insert.Close()

	return true
}

// GetMessages : Return Messages
func GetMessages(username string) []Message {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return []Message{}
	}
	defer db.Close()

	results, err := db.Query("SELECT name, sender, hash, message FROM messages WHERE name = ?", username)

	if err != nil {
		return []Message{}
	}

	var messages []Message

	for results.Next() {
		var msg Message
		err = results.Scan(&msg.To, &msg.From, &msg.Hash, &msg.Content)
		if err != nil {
			return []Message{}
		}
		messages = append(messages, msg)
	}
	return messages
}

// GetAllMessages : Return Messages from Hash
func GetAllMessages(hash string) []Message {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return []Message{}
	}
	defer db.Close()

	results, err := db.Query("SELECT name, sender, hash, message FROM messages WHERE hash = ?", hash)

	if err != nil {
		return []Message{}
	}

	var messages []Message

	for results.Next() {
		var msg Message
		err = results.Scan(&msg.To, &msg.From, &msg.Hash, &msg.Content)
		if err != nil {
			return []Message{}
		}
		messages = append(messages, msg)
	}
	return messages
}

// AddComment : Add comment to product page
func AddComment(username string, product string, content string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return false
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO comments VALUES (0, ?, ?, ?)", username, product, content)
	if err != nil {
		log.Println("insert failed")
		return false
	}
	defer insert.Close()

	return true
}

// GetComments : Get all Comments for one product
func GetComments(product string) []Comment {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return []Comment{}
	}
	defer db.Close()

	results, err := db.Query("SELECT name, product, content FROM comments WHERE product = ? ORDER BY id DESC LIMIT 20", product)

	if err != nil {
		return []Comment{}
	}

	var comments []Comment

	for results.Next() {
		var cmnt Comment
		err = results.Scan(&cmnt.User, &cmnt.Product, &cmnt.Content)
		if err != nil {
			return []Comment{}
		}
		comments = append(comments, cmnt)
	}
	return comments
}

// AddTicket : Add ticket to database
func AddTicket(username string, subject string, hash string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return false
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO tickets VALUES (?, ?, ?)", username, subject, hash)
	if err != nil {
		return false
	}
	defer insert.Close()

	return true
}

// GetTicket : Return Ticket from db if existing
func GetTicket(hash string) Ticket {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return Ticket{}
	}
	defer db.Close()

	var ticket Ticket
	err = db.QueryRow("SELECT name, subject, hash FROM tickets WHERE hash = ?", hash).Scan(&ticket.User, &ticket.Subject, &ticket.Hash)

	if err != nil {
		return Ticket{}
	}
	return ticket
}

// GetTickets : Return Tickets for user
func GetTickets(username string) []Ticket {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return []Ticket{}
	}
	defer db.Close()

	results, err := db.Query("SELECT name, subject, hash FROM tickets WHERE name = ? LIMIT 10", username)

	if err != nil {
		return []Ticket{}
	}

	var tickets []Ticket

	for results.Next() {
		var ticket Ticket
		err = results.Scan(&ticket.User, &ticket.Subject, &ticket.Hash)
		if err != nil {
			return []Ticket{}
		}
		tickets = append(tickets, ticket)
	}

	return tickets
}

// PrintDB : Print all "users" table entries
// TODO: Remove this at some point, only to be used now for debugging
func PrintDB() {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		log.Println("Error connecting to database.")
		return
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM users")

	if err != nil {
		log.Println("Query error.")
		return
	}

	for results.Next() {
		var userTest User

		err = results.Scan(&userTest.Username, &userTest.Password, &userTest.Status, &userTest.Admin)

		if err != nil {
			log.Printf("Error Scanning results.")
		}

		// log.Printf("%s %s %s %t\n", userTest.Username, userTest.Password, userTest.Status, userTest.Admin)
	}

	results, err = db.Query("SELECT * FROM messages")

	if err != nil {
		log.Println("Query error messages.")
		return
	}

	for results.Next() {
		var msg Message

		err = results.Scan(&msg.To, &msg.From, &msg.Content)

		if err != nil {
			log.Printf("Error Scanning messages results.")
		}

		log.Printf("%s %s %s\n", msg.To, msg.From, msg.Content)
	}
}
