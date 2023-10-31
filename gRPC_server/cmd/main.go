package main

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
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

func (r *Rides) Location(stream pb.Rides_LocationServer) error {
	count := int64(0)
	driverID := ""
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			err = status.Errorf(codes.Internal, "can't read")
			if err != nil {
				return err
			}
		}
		// TODO: update db
		driverID = req.DriverId
		count++
	}
	resp := pb.LocationResponse{
		DriverId: driverID,
		Count:    count,
	}
	return stream.SendAndClose(&resp)
}

type Rides struct {
	pb.UnimplementedRidesServer
}
