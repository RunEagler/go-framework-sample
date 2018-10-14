package sample

import (
	"database/sql"
	"fmt"
	"net/url"
	dockertest "github.com/ory/dockertest"

	"github.com/lib/pq"
)

//PostgresDockerPort expose port for postgres docker
const PostgresDockerPort = "5432/tcp"

// NewPostgreSQL creates and initializes a new instance of PostgreSQL.
func NewPostgreSQL(s string) (*sql.DB, error) {
	return sql.Open("postgres", s)
}

func createPostgresDocker() (*sql.DB, func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("Could not connect to docker: %s", err)
	}
	u, err := url.Parse(pool.Client.Endpoint())
	if err != nil {
		return nil, nil, fmt.Errorf("Could not parse the endpoint: %s", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_PASSWORD=passw0rd",
		"POSTGRES_USER=postgres",
		"POSTGRES_DB=postgres",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Could not start resource: %s", err)
	}
	var postgresql *sql.DB

	if err := pool.Retry(func() error {

		postgresql, err = NewPostgreSQL(PostgreSQLConnectionString(
			"postgres",
			"passw0rd",
			u.Hostname(),
			resource.GetPort(PostgresDockerPort),
			"postgres",
			"disable",
		))
		if err != nil {
			return err
		}
		return postgresql.Ping()
	}); err != nil {
		return nil, nil, fmt.Errorf("Could not connect to docker: %s", err)
	}
	return postgresql, func() error { return pool.Purge(resource) }, nil
}

// PostgreSQLConnectionString creates a connection string for PostgreSQL.
func PostgreSQLConnectionString(user, password, hostname, port, dbname, sslmode string) string {
	s, _ := pq.ParseURL(
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user,
			password,
			hostname,
			port,
			dbname,
			sslmode,
		),
	)
	return s
}
