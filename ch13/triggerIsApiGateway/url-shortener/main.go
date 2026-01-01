package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	URL string `json:"url"`
}
type Response struct {
	ShortURL string `json:"short_url"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req Request
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}
	shortURL := "https://short.ly/" + "abc123" // Placeholder logic
	res := Response{ShortURL: shortURL}
	resBody, _ := json.Marshal(res)
	return events.APIGatewayProxyResponse{
		Body:       string(resBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
