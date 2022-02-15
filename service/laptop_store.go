package service

import (
	"errors"
	"fmt"
	"github.com/jimyag/grpc-go/pb"
	"github.com/jinzhu/copier"
	"sync"
)

var ErrAlreadyExists = errors.New("record already exists")

//
// LaptopStore
//  @Description: 存储 laptop 的接口
//
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
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
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy  laptop data :%w", err)
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

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data :%w", err)
	}
	return other, nil
}
