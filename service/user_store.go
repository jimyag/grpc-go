package service

import "sync"

//
// UserStore
//  @Description:
//
type UserStore interface {
	//
	// Save
	//  @Description:
	//  @param user
	//  @return error
	//
	Save(user *User) error
	//
	// Find
	//  @Description:
	//  @param username
	//  @return *User
	//  @return error
	//
	Find(username string) (*User, error)
}

//
// InMemoryUserStore
//  @Description:
//
type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

//
// NewInmemoryUserStore
//  @Description:
//  @return *InMemoryUserStore
//
func NewInmemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

//
// Save
//  @Description:
//  @receiver store
//  @param user
//  @return error
//
func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.Username] != nil {
		return ErrAlreadyExists
	}

	store.users[user.Username] = user.Clone()
	return nil
}

//
// Find
//  @Description:
//  @receiver store
//  @param username
//  @return *User
//  @return error
//
func (store *InMemoryUserStore) Find(username string) (*User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user, ok := store.users[username]
	if !ok {
		return nil, nil
	}

	return user.Clone(), nil
}
