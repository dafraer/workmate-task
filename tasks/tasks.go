package tasks

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Constants for task status
const (
	pending  = "pending"
	running  = "running"
	finished = "finished"
)

type TaskManager struct {
	tasks map[string]*Task
	mx    sync.RWMutex
}

type Task struct {
	Id      string
	Status  string
	Payload []byte
	Result  []byte
	mx      sync.RWMutex
}

// NewTaskManager creates a new TaskManager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

// AddTask adds a new task and runs it in a seperate goroutine
func (tm *TaskManager) CreateTask(payload []byte) string {
	tm.mx.Lock()
	defer tm.mx.Unlock()
	//Create a new task with the given payload
	//Create a unique id for the task
	id := uuid.New().String()
	//Create a new task
	task := &Task{
		Id:      id,
		Status:  pending,
		Payload: payload,
	}

	//Add the task to the task manager
	tm.tasks[id] = task

	//Run the task in a separate goroutine
	go task.Run()

	//Return the id of the task
	return id
}

// Run simulates running I/O bound task
func (t *Task) Run() {
	t.mx.Lock()
	t.Status = running
	t.mx.Unlock()
	//Simulate some work
	time.Sleep(time.Minute * 3)
	t.mx.Lock()
	t.Status = finished
	t.Result = []byte("result")
	t.mx.Unlock()
}

// GetTaskResult returns the result of the task with the given id
func (tm *TaskManager) GetTaskResult(id string) (Task, error) {
	tm.mx.RLock()
	defer tm.mx.RUnlock()

	//Get the tak from the task manager
	task, ok := tm.tasks[id]

	//Check if the task exists
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}

	//Delete the task from the task manager if it is finished
	if task.Status == finished {
		tm.mx.RUnlock()
		tm.mx.Lock()
		delete(tm.tasks, id)
		tm.mx.Unlock()
	}
	task.mx.RLock()
	defer task.mx.RUnlock()
	return Task{Id: task.Id, Status: task.Status, Payload: task.Payload, Result: task.Result}, nil
}
