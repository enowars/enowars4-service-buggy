package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

type user struct {
	Username string
	Password string
	Status   string
	Admin    bool
}

// InsertUser : Insert user if not present
func InsertUser(username string, pw string, status string, admin bool) bool {
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

// AuthUser : Authenticate user
func AuthUser(username string, pw string) bool {
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

// DeleteUser : Delete user if present
func DeleteUser(username string) bool {
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

// PrintDB : Print all "users" table entries
// TODO: Remove this at some point, only to be used now for debugging
func PrintDB() {

	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/enodb")

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
		var userTest user

		err = results.Scan(&userTest.Username, &userTest.Password, &userTest.Status, &userTest.Admin)

		if err != nil {
			log.Printf("Error Scanning results.")
		}

		log.Printf("%s %s %s %t\n", userTest.Username, userTest.Password, userTest.Status, userTest.Admin)
	}
}
