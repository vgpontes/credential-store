package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

type Event struct {
}

type User struct {
	UserID   uint   `json:"userID"`
	Username string `json:"userName"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func newUser(userName string, password string) *User {
	return &User{
		UserID:   12345,
		Username: userName,
		Password: password,
		Salt:     "fewdnwqeelrdn",
	}
}

/*func handler(ctx context.Context, request Request) (Response, error) {
	fmt.Printf(request.Body)
	response := Response{StatusCode: 200, Body: "Hello World"}
	return response, nil
} */

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /authorize", func(w http.ResponseWriter, r *http.Request) {
		user := newUser("makarimi", "erhjdeer145")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})
	lambda.Start(httpadapter.New(mux).ProxyWithContext)
}
