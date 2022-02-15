package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jimyag/grpc-go/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

//
// LaptopServer
//  @Description: laptop的服务
//
type LaptopServer struct {
	Store LaptopStore
}

//
// NewLaptopServer
//  @Description: 创建一个laptopServer
//  @param store laptop的存储地方
//  @return *LaptopServer
//
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{Store: store}
}

//
// CreateLaptop
//  @Description: 创建laptop使用 req
//  @receiver service
//  @param ctx
//  @param in
//  @return *pb.CreateLaptopResponse
//  @return error
//
func (server *LaptopServer) CreateLaptop(ctx context.Context, in *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	// 从req中获得laptop信息
	laptop := in.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	// 信息校验
	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Laptop ID is not valid uuid: %v", err)
		}
	} else {
		// 没有编号
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID:%v", err)
		}
		laptop.Id = id.String()
	}

	// 测试超时
	//time.Sleep(6 * time.Second)

	// 客户端是否取消连接
	if errors.Is(ctx.Err(), context.Canceled) {
		log.Print("request is canceled")
		return nil, status.Errorf(codes.Canceled, "request is canceled")
	}

	// 判断是否超时
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Print("deadline is exceeded")
		return nil, status.Errorf(codes.DeadlineExceeded, "deadline is exceeded")
	}

	// 保存 laptop 到 store
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store :%v", err)
	}

	log.Printf("saved laptop with id : %s", laptop.Id)
	// 给予响应
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id}
	return res, nil
}

//
// SearchLaptop
//  @Description:
//  @receiver server
//  @param *pb.SearchLaptopRequest
//  @param pb.LaptopService_SearchLaptopServer
//  @return error
//
func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("recevier a search-laptop request with filter: %v", filter)

	err := server.Store.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{
			Laptop: laptop,
		}

		err := stream.Send(res)
		if err != nil {
			return err
		}

		log.Printf("sent laptop with id :%s", laptop.GetId())

		return nil
	})
	if err != nil {
		return status.Errorf(codes.Internal, "Unexpected error :%v", err)
	}
	return nil
}
