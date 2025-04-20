package tasks

import "errors"

const ExampleId = "exampleId"

type MockTaskManager struct{}

func NewMockTaskManager() *MockTaskManager {
	return &MockTaskManager{}
}

func (m *MockTaskManager) CreateTask(payload []byte) string {
	return ExampleId
}

func (m *MockTaskManager) GetTaskResult(id string) (Task, error) {
	if id != ExampleId {
		return Task{}, errors.New("id is wrong")
	}
	mockTask := Task{
		Id:      ExampleId,
		Status:  running,
		Payload: []byte("payload"),
	}
	return mockTask, nil
}
