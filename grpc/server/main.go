package main

import (
	"context"
	"log"
	"net"

	pb "go-grpc-rest-demo/grpc/user.proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	users     map[int32]*pb.User
	currentID int32
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, in *pb.UserId) (*pb.User, error) {
	user, exists := s.users[in.Id]
	if exists {
		return user, nil
	}
	return nil, grpc.Errorf(grpc.Code(grpc.NotFound), "User not found")
}

func (s *server) CreateUser(ctx context.Context, in *pb.User) (*pb.User, error) {
	s.currentID++
	in.Id = s.currentID
	s.users[s.currentID] = in
	return &pb.UserId{Id: s.currentID}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{
		users: make(map[int32]*pb.User),
	})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
