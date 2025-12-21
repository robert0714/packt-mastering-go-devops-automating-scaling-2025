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
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(4 * 1024 * 1024)),
	}
	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
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
func getOriginalURL(shortUrl string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.GetOriginalURL(
		ctx,
		&pb.GetRequest{ShortUrl: shortUrl},
	)
	if err != nil {
		log.Printf("Error fetching original URL: %v", err)
		return
	}
	log.Printf("Original URL: %s", res.OriginalUrl)
}
func main() {
	shortenURL("https://example.com")
}
