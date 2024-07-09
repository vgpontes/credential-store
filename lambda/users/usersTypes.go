package main

import (
	"time"
)

type CreateUserRequest struct {
	Username string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type GetUsersResponse struct {
	Username string `json:"userName"`
	Email    string `json:"email"`
}

type User struct {
	UserID    uint      `json:"userID"`
	Username  string    `json:"userName"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUser(username, password, email string) *User {
	return &User{
		Username:  username,
		Password:  password,
		Email:     email,
		IsAdmin:   false,
		CreatedAt: time.Now().UTC(),
	}
}
