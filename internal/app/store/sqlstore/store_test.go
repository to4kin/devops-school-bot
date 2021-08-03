package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
	migrations  string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://localhost/devops_school_test?user=postgres&password=example&sslmode=disable"
	}

	migrations = os.Getenv("TEST_MIGRATIONS")
	if migrations == "" {
		migrations = "../../../../db/migrations"
	}

	os.Exit(m.Run())
}
