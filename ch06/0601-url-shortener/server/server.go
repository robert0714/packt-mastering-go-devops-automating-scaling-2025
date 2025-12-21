package main

import (
	"context"
	"fmt"
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
