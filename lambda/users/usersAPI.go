package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIError struct {
	Error string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	db Database
}

func NewAPIServer(db Database) *APIServer {
	return &APIServer{
		db: db,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users/signup", makeHTTPHandlerFunc(s.handleCreateAccount))
	mux.HandleFunc("GET /users", makeHTTPHandlerFunc(s.handleGetUser))
	mux.HandleFunc("GET /users/{username}", makeHTTPHandlerFunc(s.handleGetUserByUsername))
	lambda.Start(httpadapter.New(mux).ProxyWithContext)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createUserReq := CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		return err
	}

	user := NewUser(createUserReq.Username, createUserReq.Password, createUserReq.Email)
	if err := s.db.CreateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	users, err := s.db.GetUsers()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetUserByUsername(w http.ResponseWriter, r *http.Request) error {
	usernameStr := r.PathValue("username")
	username, err := s.db.GetUserByUsername(usernameStr)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, username)
}
