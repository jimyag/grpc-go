package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

//
// User
//  @Description:
//
type User struct {
	Username       string
	HashedPassword string
	Role           string
}

//
// NewUser
//  @Description:
//  @param username
//  @param password
//  @param role
//  @return *User
//  @return error
//
func NewUser(username string, password string, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash passwor :%v", err)
	}
	user := &User{
		Username:       username,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}
	return user, nil
}

//
// IsCorrectPassword
//  @Description:
//  @receiver user
//  @param password
//  @return bool
//
func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

//
// Clone
//  @Description:
//  @receiver user
//  @return *User
//
func (user *User) Clone() *User {
	return &User{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
	}
}
