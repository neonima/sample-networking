package light

import (
	"sync"

	"github.com/neonima/sample-networking/pkg/store"
)

// Light hold the data a light store
type Light struct {
	// TODO: explore sync.Map + benchmark
	*sync.RWMutex
	FriendRelations map[int]map[int]*store.UserContainer
	UserInfo        map[int]*store.UserContainer
}

// New return a new Light store
func New() store.UserStore {
	return &Light{
		FriendRelations: map[int]map[int]*store.UserContainer{},
		UserInfo:        map[int]*store.UserContainer{},
		RWMutex:         &sync.RWMutex{},
	}
}

// AddUser adds a new user to the store
func (l *Light) AddUser(user *store.UserContainer) error {
	l.Lock()
	defer l.Unlock()
	l.UserInfo[user.ID] = user
	l.addRelations(user)
	return nil
}

// GetUser gets a user from the given userID from the store
func (l *Light) GetUser(userID int) (*store.UserContainer, error) {
	l.RLock()
	defer l.RUnlock()
	return l.UserInfo[userID], nil
}

// UpdateUser updates a user in the store
func (l *Light) UpdateUser(user *store.UserContainer) error {
	return nil
}

// FriendWith returns a list of friends of a given user
func (l *Light) FriendWith(user *store.UserContainer) ([]*store.UserContainer, error) {
	l.RLock()
	defer l.RUnlock()
	return l.getFriends(user), nil
}

// RemoveFriend removes a user's friend from the store
func (l *Light) RemoveFriend(user *store.UserContainer, friend int) error {
	l.Lock()
	defer l.Unlock()
	delete(l.FriendRelations[friend], user.ID)
	return nil
}

func (l *Light) addRelations(user *store.UserContainer) {
	for _, friend := range user.Friends {
		l.FriendRelations[friend] = map[int]*store.UserContainer{}
		l.FriendRelations[friend][user.ID] = user
	}
}

func (l *Light) getFriends(user *store.UserContainer) []*store.UserContainer {
	var friends []*store.UserContainer
	for _, fr := range user.Friends {
		toAdd, ok := l.FriendRelations[user.ID][fr]
		if ok && toAdd != nil {
			friends = append(friends, toAdd)
		}

	}
	return friends
}
