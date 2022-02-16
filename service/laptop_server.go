package service

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/jimyag/grpc-go/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

const maxImageSize = 1 << 30

//
// LaptopServer
//  @Description: laptop的服务
//
type LaptopServer struct {
	LaptopStore LaptopStore
	imageStore  ImageStore
	ratingStore RatingStore
}

//
// NewLaptopServer
//  @Description: 创建一个laptopServer
//  @param store laptop的存储地方
//  @return *LaptopServer
//
func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore, ratingStore RatingStore) *LaptopServer {
	return &LaptopServer{LaptopStore: laptopStore, imageStore: imageStore, ratingStore: ratingStore}
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

	if err := contextErr(ctx); err != nil {
		return nil, err
	}

	// 保存 laptop 到 store
	err := server.LaptopStore.Save(laptop)
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

	err := server.LaptopStore.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
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

//
// UploadImage
//  @Description:
//  @receiver server
//  @param stream
//  @return error
//
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logErr(status.Errorf(codes.Unknown, "cannot receive image info"))
	}
	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an image upload-image request for laptop  %s with image type %s", laptopID, imageType)

	laptop, err := server.LaptopStore.Find(laptopID)
	if err != nil {
		return logErr(status.Errorf(codes.Internal, "cannot find laptop : %v", err))
	}

	if laptop == nil {
		return logErr(status.Errorf(codes.InvalidArgument, "laptop %s doesn't exist", laptopID))
	}
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		if err := contextErr(stream.Context()); err != nil {
			return err
		}
		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}

		if err != nil {
			return logErr(status.Errorf(codes.Unknown, "cannot receive chunk data :%v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size :%d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return logErr(status.Errorf(codes.InvalidArgument, "image is too large :%d > %d", imageSize, maxImageSize))
		}

		// 测试超时
		//time.Sleep(time.Second)

		_, err = imageData.Write(chunk)
		if err != nil {
			return logErr(status.Errorf(codes.Internal, "cannot write chunk data :%v", err))
		}
	}

	imageID, err := server.imageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logErr(status.Errorf(codes.Internal, "cannot save image to the store :%v", err))
	}
	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logErr(status.Errorf(codes.Unknown, "cannot send response :%v", err))
	}

	log.Printf("saved image with id: %s , size: %d", imageID, imageSize)
	return nil
}

//
// RateLaptop
//  @Description:
//  @receiver server
//  @param server2
//  @return error
//
func (server *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		err := contextErr(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}

		if err != nil {
			return logErr(status.Errorf(codes.Unknown, "cannot receive stream request :%v", err))
		}

		laptopID := req.GetId()
		score := req.GetScore()
		log.Printf("received a rate-laptop request :id = %s ,score: = %.2f", laptopID, score)

		found, err := server.LaptopStore.Find(laptopID)
		if err != nil {
			return logErr(status.Errorf(codes.Internal, "cannot find laptop :%v", err))
		}
		if found == nil {
			return logErr(status.Errorf(codes.NotFound, "laptopID %s is not found", laptopID))
		}

		rating, err := server.ratingStore.Add(laptopID, score)
		if err != nil {
			return logErr(status.Errorf(codes.Internal, "cannot add rating to the store :%v", err))
		}

		res := &pb.RateLaptopResponse{
			LaptopId:     laptopID,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}

		err = stream.Send(res)
		if err != nil {
			return logErr(status.Errorf(codes.Internal, "cannot send stream response :%v", err))
		}

	}
	return nil
}

//
// contextErr
//  @Description:
//  @param ctx
//  @return error
//
func contextErr(ctx context.Context) error {

	switch ctx.Err() {
	case context.Canceled:
		return logErr(status.Errorf(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logErr(status.Errorf(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}

}

//
// logErr
//  @Description:
//  @param err
//  @return error
//
func logErr(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
