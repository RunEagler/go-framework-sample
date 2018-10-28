package realtime_database

import (
	"testing"
	"log"
	"context"
	"github.com/magiconair/properties/assert"
)

func TestUserRepositoryGetSet(t *testing.T) {

	type param struct {
		path string
		user User
	}

	type testCase struct {
		title string
		param param
	}

	testCases := []testCase{
		{
			title: "正常系",
			param: param{
				path: "users/1",
				user: User{
					Name: "user1",
					Age:  19,
				},
			},
		},
		{
			title: "異常系",
			param: param{
				path: "users/2",
				user: User{
					Name: "user2",
					Age:  29,
				},
			},
		},
	}
	database, err := App.Database(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, testCase := range testCases {
		userRepository := NewUserRepository(database, testCase.param.path)
		err := userRepository.Set(&testCase.param.user)
		if err != nil {
			log.Fatal(err)
		}

		actualUser, err := userRepository.Get()
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, testCase.param.user.Name, actualUser.Name, "user.name")
		assert.Equal(t, testCase.param.user.Age, actualUser.Age, "user.age")
	}

}
