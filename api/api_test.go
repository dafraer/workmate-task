package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dafraer/workmate-task/tasks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunTaskHandler(t *testing.T) {
	//Create logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	//Create a new service for testing
	s := New(logger.Sugar(), tasks.NewMockTaskManager())

	//Create test server
	server := httptest.NewServer(http.HandlerFunc(s.runTaskHandler))

	// Prepare the request payload
	payload := runRequest{[]byte("payload")}
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	// Create a POST request with JSON body
	req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	//Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	//Check that status code is OK
	assert.Equal(t, resp.StatusCode, http.StatusOK, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Decode json response into TokenPair struct
	var response runResponse
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

	//Check if tokens have been received
	assert.Equal(t, tasks.ExampleId, response.Id)

	//Close response body
	assert.NoError(t, resp.Body.Close())
}

func TestGetTaskHandler(t *testing.T) {
	//Create logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	//Create a new service for testing
	s := New(logger.Sugar(), tasks.NewMockTaskManager())

	//Create test server
	// build a ServeMux with your pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/task/get/{id}", s.getTaskHandler)
	server := httptest.NewServer(mux)

	//Send the request
	resp, err := http.Get(server.URL + "/task/get/exampleId")
	assert.NoError(t, err)

	//Check that status code is OK
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Decode json response into TokenPair struct
	var response tasks.Task
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

	//Check if tokens have been received
	assert.Equal(t, tasks.ExampleId, response.Id)

	//Close response body
	assert.NoError(t, resp.Body.Close())
}
