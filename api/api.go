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
	TaskManager *tasks.TaskManager
	logger      *zap.SugaredLogger
}

// NewService creates a new Service
func NewService(logger *zap.SugaredLogger, tm *tasks.TaskManager) *Service {
	return &Service{
		TaskManager: tm,
		logger:      logger,
	}
}

// Run starts the HTTP server
func (s *Service) Run(ctx context.Context, address string) error {
	//Create a new http server
	srv := &http.Server{
		Addr:        address,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	//Two REST routes: one for creating a task and another for getting the result of a task
	http.HandleFunc("/task/run", s.RunTaskHandler)
	http.HandleFunc("/task/get/{id}", s.GetTaskHandler)

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

// RunTaskHandler creates a new task and returns its ID
func (s *Service) RunTaskHandler(w http.ResponseWriter, r *http.Request) {
	var payload []byte
	//Decode the paylaod from the request body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Errorw("Error decoding json", "error", err)
		return
	}

	//Create a new task with the given payload
	id := s.TaskManager.CreateTask(payload)
	response, err := json.Marshal(map[string]string{"id": id})
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

// GetTaskHandler returns the result of a task by its ID
func (s *Service) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get the task ID from the URL
	id := r.URL.Query().Get("id")
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
