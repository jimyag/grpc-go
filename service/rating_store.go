package service

import "sync"

//
// RatingStore
//  @Description:
//
type RatingStore interface {
	Add(laptopID string, score float64) (*Rating, error)
}

//
// Rating
//  @Description:
//
type Rating struct {
	Count uint32
	Sum   float64
}

//
// InMemoryRatingStore
//  @Description:
//
type InMemoryRatingStore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

//
// NewInMemoryRatingStore
//  @Description:
//  @return *InMemoryRatingStore
//
func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

//
// Add
//  @Description:
//  @receiver store
//  @param laptopID
//  @param score
//  @return *Rating
//  @return error
//
func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*Rating, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	rating := store.rating[laptopID]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}

	store.rating[laptopID] = rating
	return rating, nil
}
