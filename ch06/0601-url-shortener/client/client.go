package main

import (
	"context"
	"fmt"
	"log"
	pb "url-shortener/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Client code to interact with the URL shortener service would go here.
	client := pb.NewURLShortenerClient(conn)

	shortenResponse, err := client.ShortenURL(
		context.Background(),
		&pb.ShortenRequest{OriginalUrl: "https://example.com"},
	)

	if err != nil {
		log.Fatalf("Failed to shorten URL: %v", err)
	}
	fmt.Println("Shortened URL:", shortenResponse.ShortUrl)

	getResponse, err := client.GetOriginalURL(
		context.Background(),
		&pb.GetRequest{ShortUrl: shortenResponse.ShortUrl},
	)
	if err != nil {
		log.Fatalf("Failed to get original URL: %v", err)
	}
	fmt.Println("Original URL:", getResponse.OriginalUrl)
}
