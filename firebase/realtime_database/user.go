package realtime_database

import (
	"firebase.google.com/go/db"
	"context"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserRepository struct {
	ref *db.Ref
}

func NewUserRepository(dbClient *db.Client, path string) UserRepository {

	ref := dbClient.NewRef(path)
	return UserRepository{
		ref: ref,
	}
}

func (r *UserRepository) Set(user *User) error {

	err := r.ref.Set(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepository) Get() (*User, error) {

	user := User{}
	err := r.ref.Get(context.Background(), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user map[string]interface{}) error {

	err := r.ref.Update(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
