package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/teststore"
)

func TestServer_HelloHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	srv := newServer(teststore.New())

	srv.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
