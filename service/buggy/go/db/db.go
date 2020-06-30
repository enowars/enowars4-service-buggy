package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string
	Password string
	Status   string
	Admin    bool
}

type Message struct {
	To      string
	From    string
	Hash    string
	Content string
}

type Comment struct {
	Timestamp string
	User      string
	Product   string
	Content   string
}

type Ticket struct {
	User    string
	Subject string
	Hash    string
}

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

func AuthUser(username string, pw string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return false
	}
	defer db.Close()

	var userReq User
	err = db.QueryRow("SELECT name, password FROM users WHERE BINARY name = ?", username).Scan(&userReq.Username, &userReq.Password)

	if err != nil {
		return false
	}

	if pw != userReq.Password {
		return false
	}

	return true
}

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

func AddComment(username string, product string, content string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return false
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO comments (name, product, content) VALUES (?, ?, ?)", username, product, content)
	if err != nil {
		return false
	}
	defer insert.Close()

	return true
}

func GetComments(product string) []Comment {
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(mysql:3306)/%s", os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_DATABASE")))

	if err != nil {
		return []Comment{}
	}
	defer db.Close()

	results, err := db.Query("SELECT created_at, name, product, content FROM comments WHERE product = ? ORDER BY id DESC LIMIT 100", product)

	if err != nil {
		return []Comment{}
	}

	var comments []Comment

	for results.Next() {
		var cmnt Comment
		err = results.Scan(&cmnt.Timestamp, &cmnt.User, &cmnt.Product, &cmnt.Content)
		if err != nil {
			return []Comment{}
		}
		comments = append(comments, cmnt)
	}
	return comments
}

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
