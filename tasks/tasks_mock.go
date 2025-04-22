package tasks

const ExampleID = "exampleId"

type MockTaskManager struct{}

func NewMockTaskManager() *MockTaskManager {
	return &MockTaskManager{}
}

func (m *MockTaskManager) CreateTask(payload []byte) string {
	return ExampleID
}

func (m *MockTaskManager) GetTaskResult(id string) (*TaskResult, error) {
	if id != ExampleID {
		return &TaskResult{}, ErrTaskNotFound
	}
	mockTask := &TaskResult{
		ID:      ExampleID,
		Status:  TaskStatusRunning,
		Payload: []byte("payload"),
	}
	return mockTask, nil
}

func (m *MockTaskManager) Stop() {
}
