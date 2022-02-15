package service_test

import (
	"context"
	"github.com/jimyag/grpc-go/pb"
	"github.com/jimyag/grpc-go/sample"
	"github.com/jimyag/grpc-go/serilaizer"
	"github.com/jimyag/grpc-go/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"testing"
)

//
// TestClientCreateLaptop
//  @Description:
//  @param t
//
func TestClientCreateLaptop(t *testing.T) {
	// 测试并行执行
	t.Parallel()

	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newClientLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id // 保存期待返回的 ID
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	res, err := laptopClient.CreateLaptop(context.Background(), req)

	require.NoError(t, err)              // 没有错
	require.NotNil(t, res)               // 结果不为空
	require.Equal(t, expectedID, res.Id) // 相等

	other, err := laptopServer.Store.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	requireSameLaptop(t, laptop, other)
}

//
// startTestLaptopServer
//  @Description:
//  @param t
//  @return *service.LaptopServer
//  @return string
//
func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	// 0表示 任意一个可用的端口
	listen, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	// 由于前面 listen 一定是可用的，所以这里的错误可以不用处理。
	// 这里的 server() 是阻塞的
	go grpcServer.Serve(listen)
	return laptopServer, listen.Addr().String()
}

//
// newClientLaptopClient
//  @Description: LaptopServiceClient
//  @param t
//  @param severAddress
//  @return pb.LaptopServiceClient
//
func newClientLaptopClient(t *testing.T, severAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(severAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

//
// requireSameLaptop
//  @Description: 判断两个电脑是否是同一个 这里由于 message 中还有许多序列化的字段，如果直接进行比较会出错。
// 				所以先序列化为 JSON 进行比较。还有另一个比较方法，这里没有写
//  @param t
//  @param laptop
//  @param other
//
func requireSameLaptop(t *testing.T, laptop, other *pb.Laptop) {
	laptopJSON, err := serilaizer.ProtobufToJSON(laptop)
	require.NoError(t, err)

	otherJSON, err := serilaizer.ProtobufToJSON(other)
	require.NoError(t, err)

	require.Equal(t, laptopJSON, otherJSON)
}
