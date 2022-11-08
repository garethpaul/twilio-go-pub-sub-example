package main

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

 func TestProcessGsNotification(t *testing.T) {
 	gin.SetMode(gin.TestMode)
 	r := gin.Default()
 	r.POST("/", processMessage)

 	w := httptest.NewRecorder()
 	req, _ := http.NewRequest("POST", "/", nil)
 	req.Header.Set("Content-Type", "application/json")
 	req.Header.Set("Authorization", os.Getenv("AUTH_TOKEN"))

 	// Create a fake message
 	var m PubSubMessage
 	m.Message.Data = []byte(`{"to": "+1555555555", "message": "hello-testing"}`)
 	m.Message.ID = "1234567890"

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