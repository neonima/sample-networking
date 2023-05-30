package store

import "github.com/neonima/sample-networking/pkg/model"

// UserStore is an interface that defines the API of a store
type UserStore interface {
	AddUser(user *UserContainer) error
	GetUser(userID int) (*UserContainer, error)
	UpdateUser(user *UserContainer) error
	FriendWith(user *UserContainer) ([]*UserContainer, error)
	RemoveFriend(user *UserContainer, friend int) error
}

// UserContainer merges the User model with extra information required internally
type UserContainer struct {
	*model.User
	Metadata *Metadata
}

type Metadata struct {
	ConnID string
}
