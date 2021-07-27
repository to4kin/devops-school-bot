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
	if token == "" {
		token = "1949550059:AAHTvp0Zm5ABVDKL8LVHAYkS-PEEGGZnEJE"
	}

	os.Exit(m.Run())
}
