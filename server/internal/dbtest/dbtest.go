package dbtest

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Babatunde50/book-crud/server/internal/database"
	"github.com/Babatunde50/book-crud/server/internal/docker"

	_ "github.com/lib/pq"
)

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *database.DB
	Teardown func()
	t        *testing.T
}

// StartDB spins up the Docker container and sets up the database.
func StartDB(t *testing.T) (*docker.Container, error) {
	t.Helper()
	image := "postgres:16"
	port := "5432"
	args := []string{"-e", "POSTGRES_PASSWORD=postgres"}

	c, err := docker.StartContainer(image, port, args...)
	if err != nil {
		return nil, err
	}

	t.Logf("Started container %s on port %s", c.ID, c.Host)

	return c, nil
}

// StopDB tears down the docker container
func StopDB(c *docker.Container) {
	docker.StopContainer(c.ID)
	log.Printf("Stopped container: %s", c.ID)
}

// NewTest sets up the test DB and returns a handle to it.
func NewTest(t *testing.T, c *docker.Container) *Test {
	t.Helper()

	dsn := fmt.Sprintf("postgres:postgres@%s/postgres?sslmode=disable", c.Host)
	var db *database.DB
	var err error
	for i := 0; i < 20; i++ {
		db, err = database.New(dsn, true)
		if err == nil {
			break
		}
		t.Logf("waiting for db to be ready: %v", err)
		time.Sleep(time.Second)
	}

	if err != nil {
		t.Fatalf("could not connect to database: %v", err)
	}

	t.Log("database ready")

	teardown := func() {
		db.Close()
		StopDB(c)
	}

	return &Test{
		DB:       db,
		Teardown: teardown,
		t:        t,
	}
}
