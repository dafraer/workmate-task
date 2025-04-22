package tasks

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ErrTaskNotFound is returned when a task with the specified ID does not exist
var ErrTaskNotFound = errors.New("task not found")

type TaskStatus string

// Constants for task status
const (
	TaskStatusPending  = TaskStatus("pending")
	TaskStatusRunning  = TaskStatus("running")
	TaskStatusFinished = TaskStatus("finished")
)

type TaskExecutor interface {
	CreateTask(payload []byte) string
	GetTaskResult(id string) (*TaskResult, error)
	Stop()
}

type TaskManager struct {
	tasks map[string]*task
	mx    sync.RWMutex
	wg    sync.WaitGroup
}

type task struct {
	id      string
	status  TaskStatus
	payload []byte
	result  []byte
	mx      sync.RWMutex
}

type TaskResult struct {
	ID      string
	Status  TaskStatus
	Payload []byte
	Result  []byte
}

// newTask creates a new task with the given payload
func newTask(payload []byte) *task {
	return &task{
		id:      uuid.New().String(),
		status:  TaskStatusPending,
		payload: payload,
	}
}

// SetStatus sets the status of the task
func (t *task) setStatus(status TaskStatus) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.status = status
}

// GetStatus returns the status of the task
func (t *task) getStatus() TaskStatus {
	t.mx.RLock()
	defer t.mx.RUnlock()
	return t.status
}

// GetTaskResult returns a copy of the task
func (t *task) getTaskResult() *TaskResult {
	t.mx.RLock()
	defer t.mx.RUnlock()
	return &TaskResult{ID: t.id, Status: t.status, Payload: t.payload, Result: t.result}
}

// FinishWithResult sets the status of the task to finished and sets the result from an argument
func (t *task) FinishWithResult(result []byte) {
	t.mx.Lock()
	t.status = TaskStatusFinished
	t.result = result
	t.mx.Unlock()
}

// NewTaskManager creates a new TaskManager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*task),
	}
}

// CreateTask adds a new task and runs it in a separate goroutine
func (tm *TaskManager) CreateTask(payload []byte) string {
	//Create a new task with the given payload
	//Create a unique id for the task
	//Create a new task
	task := newTask(payload)

	//Add the task to the task manager
	tm.mx.Lock()
	tm.tasks[task.id] = task
	tm.mx.Unlock()

	//Run the task in a separate goroutine
	tm.wg.Add(1)
	go func() {
		task.run()
		tm.wg.Done()
	}()

	//Return the id of the task
	return task.id
}

// Run simulates running I/O bound task
func (t *task) run() {
	t.setStatus(TaskStatusRunning)

	//Simulate some work
	time.Sleep(time.Minute * 3)

	t.FinishWithResult([]byte("result"))
}

// GetTaskResult returns the result of the task with the given id
func (tm *TaskManager) GetTaskResult(id string) (*TaskResult, error) {
	tm.mx.RLock()

	//Get the tak from the task manager
	t, ok := tm.tasks[id]
	tm.mx.RUnlock()

	//Check if the task exists
	if !ok {
		return &TaskResult{}, ErrTaskNotFound
	}

	//Delete the task from the task manager if it is finished
	if t.getStatus() == TaskStatusFinished {
		tm.mx.Lock()
		delete(tm.tasks, id)
		tm.mx.Unlock()
	}
	return t.getTaskResult(), nil
}

// Stop waits for all tasks to finish
func (tm *TaskManager) Stop() {
	tm.wg.Wait()
}
