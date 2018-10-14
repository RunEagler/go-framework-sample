package sample

import (
	"testing"
	"fmt"
	"strings"
	"io/ioutil"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
	"github.com/stretchr/testify/assert"
)

var postgresql *sql.DB
var postgresqlx *sqlx.DB
var selectBooksStmt *sqlx.NamedStmt
var selectCustomersStmt *sqlx.NamedStmt

type Book struct {
	BookID            int       `db:"book_id"`
	Title             string    `db:"title"`
	Pages             int       `db:"pages"`
	DateOfPublication time.Time `db:"date_of_publication"`
}

type Customer struct {
	CustomerID int    `db:"customer_id"`
	Name       string `db:"name"`
	Sex        string `db:"sex"`
	Age        int    `db:"age"`
}

func TestCreatePostgresDocker(t *testing.T) {

	var postgresClose func() error
	var err error
	postgresql, postgresClose, err = createPostgresDocker()
	if err != nil {
		t.Error(err)
	}

	postgresqlx = sqlx.NewDb(postgresql, "postgres")
	testCases := []struct {
		title     string
		books     []Book
		customers []Customer
	}{
		{
			title: "book tableのレコード取得",
			books: []Book{
				{
					Title:             "mind",
					Pages:             342,
					DateOfPublication: time.Date(1914, 3, 13, 0, 0, 0, 0, time.UTC),
				},
				{
					Title:             "gate",
					Pages:             685,
					DateOfPublication: time.Date(1920, 10, 10, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			title: "customer tableのレコード取得",
			customers: []Customer{
				{
					Name: "John",
					Age:  35,
					Sex:  "M",
				},
				{
					Name: "Alex",
					Age:  42,
					Sex:  "M",
				},
				{
					Name: "Jane",
					Age:  25,
					Sex:  "F",
				},
			},
		},
	}

	initTest(t)

	for _, testCase := range testCases {

		if testCase.books != nil {
			actualBooks, err := selectBooks()
			if err != nil {
				t.Error(err)
			}
			equalBook(t, testCase.title, testCase.books, actualBooks)
		} else if testCase.customers != nil {
			actualCustomers, err := selectCustomers()
			if err != nil {
				t.Error(err)
			}
			equalCustomer(t, testCase.title, testCase.customers, actualCustomers)
		}
	}

	postgresClose()
}

func equalBook(t *testing.T, title string, expect, actual []Book) {

	assert.Equal(t, len(expect), len(actual), fmt.Sprintf("[%v] :len(Book)", title))
	for i, expectBook := range expect {
		actualBook := actual[i]
		assert.Equal(t, expectBook.Title, actualBook.Title, fmt.Sprintf("[%v] :title", title))
		assert.Equal(t, expectBook.Pages, actualBook.Pages, fmt.Sprintf("[%v] :pages", title))
		assert.Equal(t, expectBook.DateOfPublication, actualBook.DateOfPublication.UTC(), fmt.Sprintf("[%v] :date_of_publication", title))
	}
}

func equalCustomer(t *testing.T, title string, expect, actual []Customer) {

	assert.Equal(t, len(expect), len(actual), fmt.Sprintf("[%v] :len(Customer)", title))
	for i, expectCustomer := range expect {
		actualCustomer := actual[i]
		assert.Equal(t, expectCustomer.Name, actualCustomer.Name, fmt.Sprintf("[%v] :name", title))
		assert.Equal(t, expectCustomer.Age, actualCustomer.Age, fmt.Sprintf("[%v] :age", title))
		assert.Equal(t, expectCustomer.Sex, actualCustomer.Sex, fmt.Sprintf("[%v] :sex", title))
	}
}

func selectBooks() ([]Book, error) {
	var books []Book
	err := selectBooksStmt.Select(&books, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return books, nil
}

func selectCustomers() ([]Customer, error) {
	var customers []Customer
	err := selectCustomersStmt.Select(&customers, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func initTest(t *testing.T) {

	var err error
	err = initTable("./tables")
	if err != nil {
		t.Error(err)
	}
	err = insertRecord("./record_yml")
	if err != nil {
		t.Error(err)
	}

	selectBooksStmt, err = postgresqlx.PrepareNamed("SELECT * FROM book")
	if err != nil {
		t.Error(err)
	}
	selectCustomersStmt, err = postgresqlx.PrepareNamed("SELECT * FROM customer")
	if err != nil {
		t.Error(err)
	}

}

func initTable(dir string) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !strings.Contains(file.Name(), "sql") {
			continue
		}
		data, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}
		_, err = postgresql.Exec(string(data))
		if err != nil {
			return err
		}
	}

	return nil

}

func insertRecord(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !strings.Contains(file.Name(), "yml") {
			continue
		}
		err = insertFixtureFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func insertFixtureFile(ymlPath string) error {

	var usedTables []string

	path := strings.Split(ymlPath, "/")
	tableName := strings.Split(path[len(path)-1], ".")[0]
	usedTables = append(usedTables, tableName)
	buf, err := ioutil.ReadFile(ymlPath)
	if err != nil {
		return err
	}
	var fixture []map[interface{}]interface{}
	err = yaml.Unmarshal(buf, &fixture)
	if err != nil {
		return err
	}
	for _, recordMap := range fixture {
		keys := []string{}
		values := []string{}
		for k, v := range recordMap {
			if v == nil {
				values = append(values, "NULL")
			} else {
				values = append(values, fmt.Sprintf("'%v'", v))
			}
			keys = append(keys, k.(string))
		}
		insertString := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", tableName, strings.Join(keys, ","), strings.Join(values, ","))
		_, err := postgresql.Exec(insertString)
		if err != nil {
			return err
		}
	}
	return nil
}
