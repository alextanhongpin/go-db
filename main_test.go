package main_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
)

var db *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	resource, err := pool.Run("postgres", "11.5-alpine", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_USER=root", "POSTGRES_DB=pg"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	if err := pool.Retry(func() error {
		log.Println("retrying")
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("host=localhost port=%s user=root password=secret dbname=pg sslmode=disable", resource.GetPort("5432/tcp")))
		// db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			log.Println("retryError:", err)
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestSomething(t *testing.T) {
	var sum int
	err := db.QueryRow("SELECT 1+1").Scan(&sum)
	if err != nil {
		t.Fatal(err)
	}
	if sum != 2 {
		t.Fatalf("want 2, got %d", sum)
	}
}
