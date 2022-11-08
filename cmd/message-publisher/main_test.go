package main

// write tests for message-publisher
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPublishReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/new-review", publishReview)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/new-review", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("AUTH_TOKEN"))

	// Create a fake message
	var m MessageContent
	m.ToNumber = "+1555555555"
	m.Message = "hello-testing"

	// Convert the message to JSON
	json, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	// Set the request body
	req.Body = ioutil.NopCloser(bytes.NewBuffer(json))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}