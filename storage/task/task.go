package task

import (
	"errors"
	"sync"
	"time"

	"github.com/bagardavidyanisntreal/clustertask/dto"
)

func NewStorage(taskCount int) *Storage {
	storage := &Storage{
		tasks: make([]*dto.Task, taskCount),
	}

	for i := range storage.tasks {
		taskDur := time.Duration(i+1) * time.Second // for diff task run freq
		storage.tasks[i] = dto.NewTask(i, taskDur)
	}

	return storage
}

type Storage struct {
	lock  sync.Mutex
	tasks []*dto.Task
}

var ErrCannotAcquire = errors.New("cannot acquire tasks storage")

func (t *Storage) Tasks() ([]*dto.Task, error) {
	if !t.lock.TryLock() {
		return nil, ErrCannotAcquire
	}

	time.Sleep(100 * time.Millisecond) // work simulation

	tasks := make([]*dto.Task, len(t.tasks))
	for i := range tasks {
		task := &t.tasks[i]
		tasks[i] = *task
	}

	t.lock.Unlock()

	return tasks, nil
}
