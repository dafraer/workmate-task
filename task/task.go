package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

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
	mx      sync.Mutex
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}
func (tm *TaskManager) NewTask(payload []byte) string {
	tm.mx.Lock()
	defer tm.mx.Unlock()
	//Create a new task with the given payload
	id := uuid.New().String()
	task := &Task{
		Id:      id,
		Status:  pending,
		Payload: payload,
	}
	tm.tasks[id] = task
	go task.Run()
	return id
}
func (t *Task) Run() {
	//Simulate I/O bound task
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

func (tm *TaskManager) GetTaskResult(id string) (Task, error) {
	tm.mx.RLock()
	defer tm.mx.RUnlock()
	task, ok := tm.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}
	if task.Status == finished {
		tm.mx.RUnlock()
		tm.mx.Lock()
		delete(tm.tasks, id)
		tm.mx.Unlock()
	}
	task.mx.Lock()
	defer task.mx.Unlock()
	return Task{Id: task.Id, Status: task.Status, Payload: task.Payload, Result: task.Result}, nil
}
