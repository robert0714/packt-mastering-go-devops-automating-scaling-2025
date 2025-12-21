package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	pb "url-shortener/proto"

	"google.golang.org/grpc"
)

type URLShortenerServer struct {
	pb.UnimplementedURLShortenerServer
	urlMap map[string]string
}

func (s *URLShortenerServer) ShortenURL(
	ctx context.Context,
	req *pb.ShortenRequest,
) (*pb.ShortenResponse, error) {
	shortURL := fmt.Sprintf("short.ly/%d", len(s.urlMap)+1)
	s.urlMap[shortURL] = req.OriginalUrl
	return &pb.ShortenResponse{ShortUrl: shortURL}, nil
}

func (s *URLShortenerServer) GetOriginalURL(
	ctx context.Context,
	req *pb.GetRequest,
) (*pb.GetResponse, error) {
	originalURL, exists := s.urlMap[req.ShortUrl]
	if !exists {
		return nil, fmt.Errorf("URL not found")
	}
	return &pb.GetResponse{OriginalUrl: originalURL}, nil
}

// ListShortenedURLs streams all shortened URLs to the client
func (s *URLShortenerServer) ListShortenedURLs(
	req *pb.ListRequest,
	stream pb.URLShortener_ListShortenedURLsServer,
) error {
	for shortURL, _ := range s.urlMap {
		resp := &pb.ShortenResponse{
			ShortUrl: shortURL,
		}
		// Send the response to the client
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

// BatchShortenURLs shortens multiple URLs
func (s *URLShortenerServer) BatchShortenURLs(
	stream pb.URLShortener_BatchShortenURLsServer,
) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&pb.BatchResponse{
					Count: int32(count),
				},
			)
		}
		if err != nil {
			return err
		}
		shortURL := fmt.Sprintf("short.ly/%d", len(s.urlMap)+1)
		s.urlMap[shortURL] = req.OriginalUrl
		count++
	}
}

// MonitorURLStats streams real-time statistics
func (s *URLShortenerServer) MonitorURLStats(
	stream pb.URLShortener_MonitorURLStatsServer,
) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		stats := &pb.StatsResponse{
			ShortUrl: req.ShortUrl,
			Clicks:   int32(len(req.ShortUrl) * 10), // Simulated clicks
			// for demo
		}
		// Send stats back to the client
		if err := stream.Send(stats); err != nil {
			return err
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterURLShortenerServer(
		server,
		&URLShortenerServer{urlMap: make(map[string]string)},
	)
	log.Println("gRPC server is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
