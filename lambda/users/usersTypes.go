package main

import (
	"time"
)

type CreateUserRequest struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type User struct {
	UserID    uint      `json:"userID"`
	Username  string    `json:"userName"`
	Password  string    `json:"password"`
	Salt      string    `json:"salt"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUser(username, password string) *User {
	return &User{
		Username:  username,
		Password:  password,
		Salt:      "dhwhdhbdfbeb",
		IsAdmin:   false,
		CreatedAt: time.Now().UTC(),
	}
}
