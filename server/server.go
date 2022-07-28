package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/karthikuppalapati/login-signup-api/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"math"
	"net"
	"time"
)

type User struct {
	Name     string
	UserId   int32
	Email    string
	Password string
}

var users []User

type Server struct {
	protobuf.LoginSignupServiceServer
}

func (s Server) SignupUser(ctx context.Context, r *protobuf.SignupUserRequest) (*protobuf.SignUpUserResponse, error) {
	fmt.Println("SignupUser triggered")
	email := r.GetUser().GetEmail()
	for _, v := range users {
		if v.Email == email {
			return nil, errors.New("Email already exists")
		}
	}
	users = append(users, User{Name: r.GetUser().GetName(), Email: r.GetUser().GetEmail(), UserId: r.GetUser().GetUserId(), Password: r.GetUser().GetPassword()})
	return &protobuf.SignUpUserResponse{Email: email}, nil
}

func (s Server) LoginUser(ctx context.Context, r *protobuf.LoginUserRequest) (*protobuf.LoginUserResponse, error) {
	fmt.Println("LoginUser triggered")
	email := r.GetEmail()
	password := r.GetPassword()

	for _, v := range users {
		if v.Email == email && v.Password == password {
			name := v.Name
			return &protobuf.LoginUserResponse{Name: name}, nil
		}
	}

	return nil, errors.New("Email or password or both are incorrect")
}

func (s Server) UpdateUserPassword(ctx context.Context, r *protobuf.UpdateUserPasswordRequest) (*protobuf.UpdateUserPasswordResponse, error) {
	fmt.Println("UpdateUserPassword triggered")

	email := r.GetEmail()
	oldPassword := r.GetOldPassword()
	newPassword := r.GetNewPassword()

	for _, v := range users {
		if v.Email == email && v.Password == oldPassword {
			v.Password = newPassword
			return &protobuf.UpdateUserPasswordResponse{Email: email}, nil
		}
	}

	return nil, errors.New("Email or password does not match")
}

func (s Server) DeleteUser(ctx context.Context, r *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	fmt.Println("DeleteUser triggered")

	email := r.GetEmail()
	password := r.GetPassword()

	for k, v := range users {
		if v.Email == email && v.Password == password {
			users = append(users[:k], users[k+1:]...)
			return &protobuf.DeleteUserResponse{Email: v.Email}, nil
		}
	}
	return nil, errors.New("Can not find the user. Email or Password is incorrect")
}

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		errors.New("Listening failed")
	}
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt64),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Timeout: 5 * time.Second,
			},
		),
	}

	s := grpc.NewServer(opts...)

	protobuf.RegisterLoginSignupServiceServer(s, &Server{})

	fmt.Println("Starting Server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
