package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jimyag/grpc-go/pb"
	"github.com/jinzhu/copier"
	"log"
	"sync"
)

var ErrAlreadyExists = errors.New("record already exists")

//
// LaptopStore
//  @Description: 存储 laptop 的接口
//
type LaptopStore interface {
	//
	// Save
	//  @Description: 保存 laptop
	//  @param laptop
	//  @return error
	//
	Save(laptop *pb.Laptop) error

	//
	// Find
	//  @Description: 查找 laptop
	//  @param id
	//  @return *pb.Laptop
	//  @return error
	//
	Find(id string) (*pb.Laptop, error)

	//
	// Search
	//  @Description: 通过过滤器搜索 laptop
	//  @param filter
	//  @param found
	//  @return error
	//
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

//
// InMemoryLaptopStore
//  @Description: 将laptop存在内存中
//
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

//
// NewInMemoryLaptopStore
//  @Description:
//  @return *InMemoryLaptopStore
//
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		mutex: sync.RWMutex{},
		data:  make(map[string]*pb.Laptop),
	}
}

//
// Save
//  @Description: 保存laptop到内存中
//  @receiver store
//  @param laptop
//  @return error
//
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// 校验数据
	if _, ok := store.data[laptop.Id]; ok {
		return ErrAlreadyExists
	}

	// 拷贝一份再保存
	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}
	store.data[other.Id] = other
	return nil
}

//
// Find
//  @Description: 查找某个laptop
//  @receiver store
//  @param id
//  @return *pb.Laptop
//  @return error
//
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop, ok := store.data[id]
	if !ok {
		return nil, nil
	}

	return deepCopy(laptop)
}

//
// Search
//  @Description: 通过过滤器搜索 laptop
//  @receiver store
//  @param filter
//  @param found
//  @return error
//
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		// 检查超时
		//time.Sleep(time.Second)
		//log.Print("check laptop id:", laptop.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is canceled")
			return fmt.Errorf("context is canceled")
		}
		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//
// isQualified
//  @Description: 是否符合过滤器的条件
//  @param filter
//  @param laptop
//  @return bool
//
func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < float32(filter.GetMinCpuGhz()) {
		return false
	}

	if toBit(laptop.GetMemory()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

//
// toBit
//  @Description: 统一单位
//  @param memory
//  @return uint64
//
func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3
	case pb.Memory_KB:
		return value << 13
	case pb.Memory_MB:
		return value << 23
	case pb.Memory_GB:
		return value << 33
	case pb.Memory_TB:
		return value << 43
	default:
		return 0

	}
}

//
// deepCopy
//  @Description:
//  @param laptop
//  @return *pb.Laptop
//  @return error
//
func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data :%w", err)
	}
	return other, nil
}
