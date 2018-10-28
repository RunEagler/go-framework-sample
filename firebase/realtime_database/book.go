package realtime_database

import (
	"time"
	"firebase.google.com/go/db"
	"context"
)

type BookRepository struct {
	ref *db.Ref
}

type Book struct {
	Title             string    `json:"title"`
	IsRead            bool      `json:"is_read"`
	DateOfPublication time.Time `json:"date_of_publication"`
	Category          Category  `json:"category"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewBookRepository(dbClient *db.Client, path string) BookRepository {

	ref := dbClient.NewRef(path)
	return BookRepository{
		ref: ref,
	}
}

func (r *BookRepository) Get() (map[string]Book, error) {

	var books map[string]Book
	err := r.ref.Get(context.Background(), &books)
	if err != nil {
		return books, err
	}
	return books, nil
}

func (r *BookRepository) Push(book *Book) error {

	_,err := r.ref.Push(context.Background(), book)
	if err != nil {
		return err
	}
	return nil
}
