package main

import (
	"context"

	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request = events.APIGatewayProxyRequest
type Response = events.APIGatewayProxyResponse

type Event struct {
}

func handler(ctx context.Context, request Request) (Response, error) {
	fmt.Printf(request.Body)
	response := Response{StatusCode: 200, Body: "Hello World"}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
