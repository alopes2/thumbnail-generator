package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type myEvent struct {
	Name string `json:"name"`
}

type response struct {
	Body       string `json:"body"`
	StatusCode int    `json:"statusCode"`
}

func handleRequest(ctx context.Context, event *myEvent) (*response, error) {
	message := response{
		Body:       "Hello from Lambda!",
		StatusCode: 200,
	}

	return &message, nil
}

func main() {
	lambda.Start(handleRequest)
}
