package main

import (
	"context"
	"log"
	"time"
	pb "url-shortener/proto"

	"google.golang.org/grpc"
)

var conn *grpc.ClientConn
var client pb.URLShortenerClient

func init() {
	var err error
	conn, err = grpc.Dial(
		"localhost:50051",
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	client = pb.NewURLShortenerClient(conn)
}
func shortenURL(url string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.ShortenURL(
		ctx,
		&pb.ShortenRequest{OriginalUrl: url},
	)
	if err != nil {
		log.Fatalf("Error while shortening URL: %v", err)
	}
	log.Printf("Short URL: %s", res.ShortUrl)
}
func main() {
	shortenURL("https://example.com")
}
