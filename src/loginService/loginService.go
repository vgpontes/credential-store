package loginservice

import "github.com/aws/aws-lambda-go/lambda"

type Event struct {
}

func handler(event *Event) {

}

func main() {
	lambda.Start(handler)
}
