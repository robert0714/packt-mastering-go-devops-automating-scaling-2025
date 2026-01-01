package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("New object uploaded: %s\n", s3.Object.Key)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
