package main

import (
	"context"
	"fmt"
	"github.com/karthikuppalapati/login-signup-api/protobuf"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Starting Client\n")

	cc, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	// Close connection before exiting application
	defer cc.Close()

	c := protobuf.NewLoginSignupServiceClient(cc)

	user1 := protobuf.User{Name: "karthik", Email: "karthiku19@gmail.com", UserId: 1, Password: "password123"}
	user2 := protobuf.User{Name: "atharv", Email: "atharva123@gmail.com", UserId: 2, Password: "athdon"}

	SignupUserResponse1, err := c.SignupUser(context.Background(), &protobuf.SignupUserRequest{User: &user1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Signing up successful with email : ", SignupUserResponse1.Email)

	SignupUserResponse2, err := c.SignupUser(context.Background(), &protobuf.SignupUserRequest{User: &user2})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Signing up successful with email : ", SignupUserResponse2.Email)

	LoginResponse, err := c.LoginUser(context.Background(), &protobuf.LoginUserRequest{Email: "karthiku19@gmail.com", Password: "password123"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Login successful. Welcome ", LoginResponse.Name)

	UpdatePasswordResponse, err := c.UpdateUserPassword(context.Background(), &protobuf.UpdateUserPasswordRequest{
		Email:       "atharva123@gmail.com",
		OldPassword: "athdon",
		NewPassword: "athdon123",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Update password successfully for : ", UpdatePasswordResponse.Email)

	DeleteUserResponse, err := c.DeleteUser(context.Background(), &protobuf.DeleteUserRequest{
		Email:    "karthiku19@gmail.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted user with email : ", DeleteUserResponse.Email)

}
