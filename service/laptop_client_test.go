package service_test

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jimyag/grpc-go/pb"
	"github.com/jimyag/grpc-go/sample"
	"github.com/jimyag/grpc-go/serilaizer"
	"github.com/jimyag/grpc-go/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"path/filepath"
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
	laptopStore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil)
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

	other, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	requireSameLaptop(t, laptop, other)
}

//
// TestLaptopServer_SearchLaptop
//  @Description: 测试有两种方法，第一种是 mock pb/laptop_service.pb.go 中 LaptopService_SearchLaptopServer 的流接口。
// 				还需要添加此接口的实现
//	type ServerStream interface {
//    SetHeader(metadata.MD) error
//    SendHeader(metadata.MD) error
//    SetTrailer(metadata.MD)
//    Context() context.Context
//    SendMsg(m interface{}) error
//    RecvMsg(m interface{}) error
//}
//	第二种方法 使用客户端调用测试服务器的RPC
//  @param t
//
func TestLaptopServer_SearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GB,
		},
	}

	store := service.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)

	laptopCount := 6
	for i := 0; i < laptopCount; i++ {
		laptop := sample.NewLaptop()
		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Memory = &pb.Memory{
				Value: 4096,
				Unit:  pb.Memory_MB,
			}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberThreads = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Memory = &pb.Memory{
				Value: 16,
				Unit:  pb.Memory_GB,
			}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberThreads = 5
			laptop.Cpu.MinGhz = 2.7
			laptop.Cpu.MaxGhz = 4.9
			laptop.Memory = &pb.Memory{
				Value: 8,
				Unit:  pb.Memory_GB,
			}
			expectedIDs[laptop.Id] = true

		}
		err := store.Save(laptop)
		require.NoError(t, err)
	}

	serverAddress := startTestLaptopServer(t, store, nil)
	laptopClient := newClientLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())
		found += 1
	}

	require.Equal(t, len(expectedIDs), found)
}

func TestClientUploadImage(t *testing.T) {
	t.Parallel()
	testImageFolder := "../tmp"

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore(testImageFolder)

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, imageStore)
	laptopClient := newClientLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/laptop.jpg", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)
	req := &pb.UploadImageRequest{Data: &pb.UploadImageRequest_Info{
		Info: &pb.ImageInfo{
			LaptopId:  laptop.GetId(),
			ImageType: imageType,
		},
	}}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024) // 1KB

	size := 0
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)

	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.Equal(t, size, int(res.GetSize()))

	saveImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, saveImagePath)
	//require.NoError(t, os.Remove(saveImagePath))
}

//
// startTestLaptopServer
//  @Description:
//  @param t
//  @return *service.LaptopServer
//  @return string
//
func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore, imageStore service.ImageStore) string {
	laptopServer := service.NewLaptopServer(laptopStore, imageStore)
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	// 0表示 任意一个可用的端口
	listen, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	// 由于前面 listen 一定是可用的，所以这里的错误可以不用处理。
	// 这里的 server() 是阻塞的
	go grpcServer.Serve(listen)
	return listen.Addr().String()
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
