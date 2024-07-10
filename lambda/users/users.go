package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := NewUsersTable()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	server := NewAPIServer(db)
	server.Run()
}

func (s *UserAPIServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", makeHTTPHandlerFunc(s.handlePutUser))
	mux.HandleFunc("GET /users", makeHTTPHandlerFunc(s.handleGetAllUsers))
	mux.HandleFunc("GET /users/{username}", makeHTTPHandlerFunc(s.handleGetUser))
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

func (s *UserAPIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	usernameStr := r.PathValue("username")
	user, err := s.db.GetUser(usernameStr)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, user)
}

func (s *UserAPIServer) handleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.db.GetAllUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
}

func (s *UserAPIServer) handlePutUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		return err
	}

	isPayloadValid, invalidPayloadResponse := validatePayload(createUserReq)

	if !isPayloadValid {
		return WriteJSON(w, http.StatusBadRequest, invalidPayloadResponse)
	}

	hashedPassword, err := saltAndHash(createUserReq.Password)
	if err != nil {
		return err
	}

	user := NewUser(createUserReq.Username, hashedPassword, createUserReq.Email)
	if err := s.db.PutUser(user); err != nil {
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
