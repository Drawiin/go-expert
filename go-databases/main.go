package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/godb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	user := User{}.NewUser("test", "test@emai.com", "password")
	if err := insertUser(db, user) ; err != nil {
		panic(err)
	}
}

type User struct{ 
	id string
	name string
	email string
	password string
	createAt string
}

func (u User) NewUser(name, email, password string) *User{
	return &User{ 
		id: uuid.New().String(),
		name: name,
		email: email,
		password: hex.EncodeToString(sha256.New().Sum([]byte(password))),
		createAt: time.Now().String(),
	}
}

func insertUser(db *sql.DB, user *User) error {
	if err := createTable(db); err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO users(id, name, email, password, createAt) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.id, user.name, user.email, user.password, user.createAt)
	if err != nil {
		return err
	}
	return nil
}

func createTable(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users(id VARCHAR(36) PRIMARY KEY, name VARCHAR(255), email VARCHAR(255), password VARCHAR(255), createAt VARCHAR(255))")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}