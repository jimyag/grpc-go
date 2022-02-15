package service

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
)

//
// ImageStore
//  @Description: store laptop image
//
type ImageStore interface {
	//
	// Save
	//  @Description: save laptop to store
	//  @param laptopID
	//  @param imageType
	//  @param imageData
	//  @return string
	//  @return error
	//
	Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error)
}

type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

//
// ImageInfo
//  @Description: the information laptop images
//
type ImageInfo struct {
	LaptopID string
	Type     string
	Path     string
}

//
// NewDiskImageStore
//  @Description: a new DiskImageStore
//  @param imageFolder
//  @return *DiskImageStore
//
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

//
// Save
//  @Description:
//  @receiver store
//  @param laptopID
//  @param imageType
//  @param imageData
//  @return string
//  @return error
//
func (store *DiskImageStore) Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generator image id:%w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create file :%w", err)
	}
	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file :%w", err)
	}
	store.mutex.RLock()
	defer store.mutex.RUnlock()
	store.images[imageID.String()] = &ImageInfo{
		LaptopID: laptopID,
		Type:     imageType,
		Path:     imagePath,
	}

	return imageID.String(), nil
}
