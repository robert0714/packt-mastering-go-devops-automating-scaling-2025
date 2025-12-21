package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
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

func listShortenedURLs(client pb.URLShortenerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.ListShortenedURLs(ctx, &pb.ListRequest{})
	if err != nil {
		log.Fatalf("Failed to list URLs: %v", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}
		fmt.Println("Shortened URL:", resp.ShortUrl)
	}
}
func batchShortenURLs(client pb.URLShortenerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.BatchShortenURLs(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	urls := []string{
		"https://example1.com",
		"https://example2.com",
		"https://example3.com",
	}
	for _, url := range urls {
		req := &pb.ShortenRequest{OriginalUrl: url}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send URL: %v", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	fmt.Printf(
		"Batch shortening completed. %d URLs shortened.\n",
		resp.Count,
	)
}
func monitorURLStats(client pb.URLShortenerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.MonitorURLStats(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	urls := []string{
		"short.ly/1",
		"short.ly/2",
	}
	// Send URLs and receive stats
	waitc := make(chan struct{})
	go func() {
		for _, url := range urls {
			req := &pb.StatsRequest{ShortUrl: url}
			if err := stream.Send(req); err != nil {
				log.Fatalf("Failed to send stats request: %v", err)
			}
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving stats: %v", err)
			}
			fmt.Printf("Stats for %s: %d clicks\n", resp.ShortUrl, resp.Clicks)
		}
		close(waitc)
	}()
	<-waitc
}
