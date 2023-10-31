package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"

	"google.golang.org/grpc"

	"gRPC/gRPC_server/pb"
)

func main() {
	addr := ":9292"

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("error: can't listen - %s", err)
	}

	srv := grpc.NewServer()
	var u Rides
	pb.RegisterRidesServer(srv, &u)
	reflection.Register(srv)

	log.Printf("info: server ready on %s", addr)
	if err = srv.Serve(lis); err != nil {
		log.Fatalf("error: can't serve - %s", err)
	}
}
func (r *Rides) Start(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata")
	}
	log.Printf("info:key is %s", md["key"])
	// TODO: Validate req
	resp := pb.StartResponse{
		Id: req.Id,
	}

	// TODO: Work (insert to database ...)
	return &resp, nil
}

type Rides struct {
	pb.UnimplementedRidesServer
}
