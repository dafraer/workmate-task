package tasks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTask(t *testing.T) {
	tm := NewTaskManager()
	payload := []byte("payload")
	//Test task creation
	id := tm.CreateTask(payload)
	assert.NoError(t, uuid.Validate(id))

	//Test task result
	task, err := tm.GetTaskResult(id)
	assert.NoError(t, err)
	assert.Equal(t, id, task.Id)
	assert.NotEmpty(t, task.Status)
	assert.Equal(t, payload, task.Payload)
	assert.Nil(t, task.Result)
}
