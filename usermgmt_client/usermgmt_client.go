package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var newUsers = make(map[string]int32)
	newUsers["John"] = 54
	newUsers["Veronica"] = 23
	for name, age := range newUsers {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf(`User Details: 
NAME: %s
AGE: %d
ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}
	// 7. Before calling service method we first initialise an empty message as input to call to get users.
	params := &pb.GetUsersParams{}
	// And then immediately following that we call GetUsers() service method.
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not retreive users: %v", err)
	}
	// We simply print the users list from the server
	log.Print("\nUSER LIST: \n")
	fmt.Printf("r.GetUsers(): %v\n", r.GetUsers())
}

// 8. After all we run the server in terminal:
// go run usermgmt_server/usermgmt_server.go

// And after that we run the client in other terminal:
// go run usermgmt_client/usermgmt_client.go
