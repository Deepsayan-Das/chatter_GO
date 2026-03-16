package models

import "time"

//sql query to create table
// CREATE TABLE users (
//     id SERIAL PRIMARY KEY,
//     username TEXT UNIQUE NOT NULL,
//     email TEXT UNIQUE NOT NULL,
//     password_hash TEXT NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
