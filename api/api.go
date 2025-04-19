package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dafraer/workmate-task/task"
	"go.uber.org/zap"
)

type Service struct {
	TaskManager *task.TaskManager
	logger      *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{
		TaskManager: task.NewTaskManager(),
		logger:      logger,
	}
}

func (s *Service) Run(ctx context.Context, addres string) error {
	http.HandleFunc("/task/run", s.RunTaskHandler)
	http.HandleFunc("/task/get/{id}", s.GetTaskHandler)
	if err := http.ListenAndServe(addres, nil); err != nil {
		return err
	}
	return nil
}

func (s *Service) RunTaskHandler(w http.ResponseWriter, r *http.Request) {
	var payload []byte
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Errorw("Error decoding json", "error", err)
		return
	}
	id := s.TaskManager.NewTask(payload)
	response, err := json.Marshal(map[string]string{"id": id})
	if err != nil {
		http.Error(w, "Error marshaling json", http.StatusInternalServerError)
		return
	}

	//Write response
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		s.logger.Errorw("Error writing a response", "error", err)
	}
}

func (s *Service) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	task, err := s.TaskManager.GetTaskResult(id)
	if err != nil {
		http.Error(w, "Error getting task", http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Error marshaling json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(response); err != nil {
		s.logger.Errorw("Error writing a response", "error", err)
	}
}
