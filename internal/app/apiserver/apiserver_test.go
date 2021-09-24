package apiserver

import (
	"os"
	"testing"
)

var (
	token string
)

func TestMain(m *testing.M) {
	token = os.Getenv("TEST_TELEGRAM_TOKEN")

	os.Exit(m.Run())
}
