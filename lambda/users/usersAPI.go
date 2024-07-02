package users

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
	mux.HandleFunc("POST /account/signup", makeHTTPHandlerFunc(s.handleCreateAccount))
	lambda.Start(httpadapter.New(mux).ProxyWithContext)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
