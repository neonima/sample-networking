package light_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/neonima/sample-networking/pkg/model"

	"github.com/neonima/sample-networking/pkg/store/light"
	"github.com/stretchr/testify/require"

	"github.com/neonima/sample-networking/pkg/store"
)

func newUser(t *testing.T, online bool) *store.UserContainer {
	t.Helper()
	return &store.UserContainer{
		User: &model.User{
			ID:      1,
			Friends: []int{45, 6, 7, 3},
			Online:  online,
		},
		Metadata: &store.Metadata{
			ConnID: "1",
		},
	}
}

func TestNew(t *testing.T) {
	expected := &light.Light{
		FriendRelations: map[int]map[int]*store.UserContainer{},
		UserInfo:        map[int]*store.UserContainer{},
		RWMutex:         &sync.RWMutex{},
	}
	received := light.New()

	require.Equal(t, expected, received)
}

func TestLight(t *testing.T) {
	t.Run("AddUser", func(t *testing.T) {
		t.Run("should add user successfully", func(t *testing.T) {
			expected := newUser(t, false)
			st := light.New()
			require.NoError(t, st.AddUser(expected))
			stL := st.(*light.Light)
			require.Equal(t, expected, stL.UserInfo[1])
			require.Len(t, stL.FriendRelations, len(expected.Friends))
		})
	})

	t.Run("GetUser", func(t *testing.T) {
		t.Run("should get a user successfully", func(t *testing.T) {
			expected := newUser(t, false)
			st := light.New()
			stL := st.(*light.Light)
			stL.UserInfo[1] = expected
			user, err := st.GetUser(1)
			require.NoError(t, err)
			require.Equal(t, expected, user)
		})
	})

	t.Run("UpdateUser", func(t *testing.T) {
		t.Run("should update a user successfully", func(t *testing.T) {
			expected := newUser(t, false)
			st := light.New()
			stL := st.(*light.Light)
			stL.UserInfo[1] = expected
			expected.Online = true
			require.NoError(t, st.UpdateUser(expected))
			require.Equal(t, expected, stL.UserInfo[1])
		})
	})

	t.Run("FriendWith", func(t *testing.T) {
		t.Run("should return an empty list of friends successfully", func(t *testing.T) {
			expected := newUser(t, false)
			st := light.New()
			require.NoError(t, st.AddUser(expected))
			friends, err := st.FriendWith(expected)
			require.NoError(t, err)
			require.Equal(t, []*store.UserContainer(nil), friends)
		})
		t.Run("should return a list of friends successfully", func(t *testing.T) {
			expected := newUser(t, false)
			st := light.New()
			stL := st.(*light.Light)
			require.NoError(t, st.AddUser(expected))

			rel := map[int]map[int]*store.UserContainer{}
			rel[1] = map[int]*store.UserContainer{}
			var expectedFriends []*store.UserContainer
			for _, fr := range expected.Friends {
				friend := &store.UserContainer{
					User: &model.User{
						ID:      fr,
						Friends: []int{1},
						Online:  true,
					},
					Metadata: &store.Metadata{ConnID: fmt.Sprintf("%v", fr)},
				}
				expectedFriends = append(expectedFriends, friend)
				rel[1][fr] = friend
			}
			stL.FriendRelations = rel

			friends, err := st.FriendWith(expected)
			require.NoError(t, err)
			require.Equal(t, expectedFriends, friends)
		})
	})

	t.Run("RemoveFriend", func(t *testing.T) {
		t.Run("should remove a user's friend successfully", func(t *testing.T) {
			t.Skip("not implemented")
		})
	})

	t.Run("SetOnline", func(t *testing.T) {
		t.Run("should set a user online successfully", func(t *testing.T) {
			t.Skip("not implemented")
		})
	})
}
