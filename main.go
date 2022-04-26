package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pgdevelopers/ly-gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

type server struct {
	sessionv1.UnimplementedSessionServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) GetSession(ctx context.Context, in *sessionv1.GetSessionRequest) (*sessionv1.GetSessionResponse, error) {
	return &sessionv1.GetSessionResponse{Session: &sessionv1.Session{
		DeviceId:         "",
		DeviceType:       "",
		ReceivedTime:     nil,
		ThingName:        "",
		TraceId:          "",
		ConsumerId:       "",
		SessionId:        "",
		Client:           "",
		ClientVersion:    "",
		Message:          "",
		SessionStartTime: nil,
		SessionType:      "",
	}}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	sessionv1.RegisterSessionServiceServer(s, &server{})
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = sessionv1.RegisterSessionServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
