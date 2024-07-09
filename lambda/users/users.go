package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", makeHTTPHandlerFunc(HandlePutUser))
	mux.HandleFunc("GET /users", makeHTTPHandlerFunc(HandleGetAllUsers))
	mux.HandleFunc("GET /users/{username}", makeHTTPHandlerFunc(HandleGetUser))
	lambda.Start(httpadapter.New(mux).ProxyWithContext)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, "Unknown exception")
		}
	}
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	db, err := NewUsersTable()
	if err != nil {
		return err
	}

	usernameStr := r.PathValue("username")
	username, err := db.GetUser(usernameStr)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, username)
}

func HandleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	db, err := NewUsersTable()
	if err != nil {
		return err
	}

	users, err := db.GetAllUsers()

	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
}

func HandlePutUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		return err
	}

	isPayloadValid, invalidPayloadResponse := validatePayload(createUserReq)

	if !isPayloadValid {
		return WriteJSON(w, http.StatusBadRequest, invalidPayloadResponse)
	}

	db, err := NewUsersTable()

	if err != nil {
		return err
	}

	hashedPassword, err := saltAndHash(createUserReq.Password)
	if err != nil {
		return err
	}

	user := NewUser(createUserReq.Username, hashedPassword, createUserReq.Email)
	if err := db.PutUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func saltAndHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func validatePayload(payload CreateUserRequest) (bool, string) {
	var invalidParams []string

	if payload.Username == "" {
		invalidParams = append(invalidParams, "username")
	}

	if payload.Password == "" {
		invalidParams = append(invalidParams, "password")
	}

	if payload.Email == "" {
		invalidParams = append(invalidParams, "email")
	}

	if len(invalidParams) > 0 {
		var invalidPayloadResponse string = "User not inserted, missing required attributes: " + strings.Join(invalidParams, ", ")
		return false, invalidPayloadResponse
	}
	return true, ""
}
