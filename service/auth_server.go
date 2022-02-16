package service

import (
	"context"
	"github.com/jimyag/grpc-go/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//
// AuthServer
//  @Description:
//
type AuthServer struct {
	userStore  UserStore
	jwtManager *JWTManager
}

//
// NewAuthServer
//  @Description:
//  @param userStore
//  @param jwtManager
//  @return *AuthServer
//
func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		jwtManager: jwtManager,
	}
}

//
// Login
//  @Description:
//  @receiver server
//  @param ctx
//  @param req
//  @return *pb.LoginResponse
//  @return error
//
func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password: %v", err)
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token: %v", err)
	}

	res := &pb.LoginResponse{AccessToken: token}
	return res, nil
}
