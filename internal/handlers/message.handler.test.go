package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSendMessageEndpoint(t *testing.T) {

	gin.SetMode(gin.TestMode)

	router := gin.Default()

	router.POST("/messages", SendMessage)

	body := []byte(`{
		"room_id":1,
		"content":"hello test"
	}`)

	req, _ := http.NewRequest(
		"POST",
		"/messages",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}
}
