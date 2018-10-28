package realtime_database

import (
	"testing"
	"context"
	"log"
	"time"
	"github.com/magiconair/properties/assert"
)

func TestBookRepositoryGetPush(t *testing.T) {

	type param struct {
		path  string
		books []Book
	}

	type testCase struct {
		title string
		param param
	}

	testCases := []testCase{
		{
			title: "正常系",
			param: param{
				path: "books/1",
				books: []Book{
					{
						Title:             "bookA",
						IsRead:            true,
						DateOfPublication: time.Date(2010, 5, 4, 0, 0, 0, 0, time.UTC),
						Category: Category{
							ID:   1,
							Name: "IT",
						},
					},
					{
						Title:             "bookB",
						IsRead:            false,
						DateOfPublication: time.Date(2004, 3, 20, 0, 0, 0, 0, time.UTC),
						Category: Category{
							ID:   2,
							Name: "Biology",
						},
					},
				},
			},
		},
	}
	database, err := App.Database(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, testCase := range testCases {
		bookRepository := NewBookRepository(database, testCase.param.path)

		for _, book := range testCase.param.books {
			err := bookRepository.Push(&book)
			if err != nil {
				log.Fatal(err)
			}

		}
		actualBooks, err := bookRepository.Get()
		if err != nil {
			log.Fatal(err)
		}
		i := 0

		assert.Equal(t, len(testCase.param.books), len(actualBooks), "len(books)")
		for _, actualBook := range actualBooks {
			expectBook := testCase.param.books[i]
			assert.Equal(t, expectBook.Title, actualBook.Title)
			assert.Equal(t, expectBook.IsRead, actualBook.IsRead)
			assert.Equal(t, expectBook.DateOfPublication, actualBook.DateOfPublication)
			assert.Equal(t, expectBook.Category.ID, actualBook.Category.ID)
			assert.Equal(t, expectBook.Category.Name, actualBook.Category.Name)
			i++
		}
	}

}
