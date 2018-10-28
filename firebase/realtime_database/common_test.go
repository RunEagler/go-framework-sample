package realtime_database

import (
	"testing"
	"google.golang.org/api/option"
	"firebase.google.com/go"
	"log"
	"context"
	"os"
)

var App *firebase.App

func TestMain(m *testing.M) {

	var err error

	opt := option.WithCredentialsFile("$SERIVCE_ACCOUNT_KEY")
	config := &firebase.Config{
		DatabaseURL: "$DATABASE_URL",
	}
	App, err = firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	result := m.Run()

	os.Exit(result)
}
