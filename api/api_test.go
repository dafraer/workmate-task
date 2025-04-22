package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dafraer/workmate-task/tasks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Decode json response into response struct
	var response runResponse
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

	//Check if the response ID matches the expected ID
	assert.Equal(t, tasks.ExampleID, response.ID)

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
	server := httptest.NewServer(http.HandlerFunc(s.getTaskHandler))

	//Send the request
	resp, err := http.Get(server.URL + "/task/get?id=" + tasks.ExampleID)
	assert.NoError(t, err)

	//Check that status code is OK
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Decode json response into task struct
	var response tasks.TaskResult
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

	//Check that task ID matches the expected ID
	assert.Equal(t, tasks.ExampleID, response.ID)

	//Check that empty request results in an status code of 400
	resp, err = http.Get(server.URL + "/task/get")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, fmt.Sprintf("expected 400 but got %d", resp.StatusCode))

	//Check that request with a non-existing id  results in an status code of 400
	resp, err = http.Get(server.URL + "/task/get?id=nonExistingId")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, fmt.Sprintf("expected 404 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())
}

// TestService tests the service by creating a real task and checking its status until it finishes
func TestService(t *testing.T) {
	//Create logger
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	//Create a new service for testing
	s := New(logger.Sugar(), tasks.NewTaskManager())

	//Create test server for run handler
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
	assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

	//Decode json response into response struct
	var response runResponse
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

	id := response.ID

	// Create test server for get handler
	server = httptest.NewServer(http.HandlerFunc(s.getTaskHandler))

	//Check that task runs properly
	for {
		//Send the request
		resp, err := http.Get(server.URL + "/task/get?id=" + id)
		assert.NoError(t, err)

		//Check that status code is OK
		assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected 200 but got %d", resp.StatusCode))

		//Decode json response into task struct
		var response tasks.TaskResult
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&response))

		//Check that taskResult matches the expected Result
		assert.Equal(t, id, response.ID)
		assert.Equal(t, []byte("payload"), response.Payload)
		if response.Status == tasks.TaskStatusFinished {
			assert.Equal(t, []byte("result"), response.Result)
			break
		}
		if response.Status == tasks.TaskStatusRunning {
			assert.Nil(t, response.Result)
		}
	}
	//Check that task is deleted after it is finished and recieved
	//Make a reqquest
	resp, err = http.Get(server.URL + "/task/get?id=" + id)
	assert.NoError(t, err)

	//Check that status code is NotFound
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, fmt.Sprintf("expected 404 but got %d", resp.StatusCode))

	//Close response body
	assert.NoError(t, resp.Body.Close())
}
