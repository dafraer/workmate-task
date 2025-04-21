package api

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/dafraer/workmate-task/tasks"
	"go.uber.org/zap"
)

type Service struct {
	TaskManager tasks.TaskExecutor
	logger      *zap.SugaredLogger
}

type runRequest struct {
	Payload []byte
}

type runResponse struct {
	Id string
}

// New creates a new Service
func New(logger *zap.SugaredLogger, tm tasks.TaskExecutor) *Service {
	return &Service{
		TaskManager: tm,
		logger:      logger,
	}
}

// Run starts an HTTP server
func (s *Service) Run(ctx context.Context, address string) error {
	//Create a new http server
	srv := &http.Server{
		Addr:        address,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	//Two REST routes: one for creating a task and another for getting the result of a task
	http.HandleFunc("/task/run", s.runTaskHandler)
	http.HandleFunc("/task/get/{id}", s.getTaskHandler)

	//Create a channel to listen for errors
	ch := make(chan error)

	//Run the server in a separate goroutine
	go func() {
		defer close(ch)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
			return
		}
		ch <- nil
	}()

	//Wait for the context to be done or for an error to occur and shutdown the server
	select {
	case <-ctx.Done():
		if err := srv.Shutdown(context.Background()); err != nil {
			return err
		}
		err := <-ch
		if err != nil {
			return err
		}
	case err := <-ch:
		return err
	}

	return nil
}

// runTaskHandler creates a new task and returns its ID
func (s *Service) runTaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskPayload runRequest
	//Decode the payload from the request body
	if err := json.NewDecoder(r.Body).Decode(&taskPayload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Errorw("Error decoding json", "error", err)
		return
	}

	//Create a new task with the given payload
	id := s.TaskManager.CreateTask(taskPayload.Payload)
	response, err := json.Marshal(runResponse{id})
	if err != nil {
		http.Error(w, "Error marshaling json", http.StatusInternalServerError)
		s.logger.Errorw("Error marshaling json", "error", err)
		return
	}

	//Write response
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		s.logger.Errorw("Error writing a response", "error", err)
	}
}

// getTaskHandler returns the result of a task by its ID
func (s *Service) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get the task ID from the URL
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// Get the task result from the TaskManager
	task, err := s.TaskManager.GetTaskResult(id)
	if err != nil {
		http.Error(w, "task with the specified ID does not exist", http.StatusInternalServerError)
		return
	}

	// Marshal the task result to JSON
	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Error marshaling json", http.StatusInternalServerError)
		s.logger.Errorw("Error marshaling json", "error", err)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		s.logger.Errorw("Error writing a response", "error", err)
	}
}
