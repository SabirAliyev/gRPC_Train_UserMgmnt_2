package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	// 2. Define a UserList attribute.
	// In this function is essentially going to be a constructor for the UserManagementServer type.

	return &UserManagementServer{
		userList: &pb.UserList{},
	}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer

	// 1. We are going to add new attribute to this type. We declare the new variable as protobuf UserList,
	// that was automatically generated using protobuf compiler. And if you hover over this type (UserList)
	// you will see, that it is an array of User struct pointers:
	// Users []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`.

	// So this User_list variable will be in memory data structure where we use to store users that we created
	// when a client called CreateNewUser() gRPC service method.
	userList *pb.UserList
}

func (server *UserManagementServer) Run() error {
	// 3. We move the main() function code into this function.

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	// And we just return the error.
	return s.Serve(lis)
}

// CreateNewUser 5. And now there is only two thins we are going to do here. The first we need to create a NewUser function.
// So, when NewUser function is called we should append a new user to UserManagement Server user list.
func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var userId = int32(rand.Intn(1000))

	// We put user definition to a new variable.
	createdUser := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: userId}

	// And we append new user to the UserManagementSerer user list and return new user.
	server.userList.Users = append(server.userList.Users, createdUser)
	return createdUser, nil
}

// GetUsers 6. The last thing we need to create getUsers function.
func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return server.userList, nil
}

func main() {
	// 4. We instantiate a new UserManagement server where we call the constructor of UserManagementServer
	var userMgmtServer = NewUserManagementServer()
	if err := userMgmtServer.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
