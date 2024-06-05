package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
}

func handler(ctx context.Context, event *Event) (*string, error) {
	print("Hello World")
	message := "Hello World"
	return &message, nil
}

func main() {
	lambda.Start(handler)
}
