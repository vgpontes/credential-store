package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type User struct {
	UserID   uint   `json:"userID"`
	Username string `json:"userName"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	IsAdmin  bool   `json:"isAdmin"`
}

var Users = []User{
	{UserID: 12345, Username: "makarimi", Password: "h2Hax.", Salt: "w3dfwdf", IsAdmin: false},
	{UserID: 12346, Username: "vpontes", Password: "d32r2r3", Salt: "qdwqwfwefwf", IsAdmin: true},
	{UserID: 12347, Username: "sakulka", Password: "dw3edwefwe", Salt: "edwedfw", IsAdmin: false},
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("username")
	for _, user := range Users {
		if user.Username == userName {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
			break
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}
	if user.IsAdmin {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Users)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{username}", handleGetUser)
	mux.HandleFunc("GET /users", handleGetUsers)
	lambda.Start(httpadapter.New(mux).ProxyWithContext)
}