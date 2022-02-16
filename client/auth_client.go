package client

import (
	"github.com/jimyag/grpc-go/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

//
// AuthClient
//  @Description:
//
type AuthClient struct {
	service  pb.AuthServiceClient
	username string
	password string
}

//
// NewAuthClient
//  @Description:
//  @param cc
//  @param username
//  @param password
//  @return *AuthClient
//
func NewAuthClient(cc *grpc.ClientConn, username string, password string) *AuthClient {
	service := pb.NewAuthServiceClient(cc)
	return &AuthClient{
		service:  service,
		username: username,
		password: password,
	}
}

//
// Login
//  @Description:
//  @receiver client
//  @return string
//  @return error
//
func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}
