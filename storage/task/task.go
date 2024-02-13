package task

import (
	"errors"
	"sync"
	"time"

	"github.com/bagardavidyanisntreal/clustertask/task"
)

func NewStorage(taskCount int) (*Storage, time.Duration) {
	storage := &Storage{
		tasks: make([]*task.Task, taskCount),
	}

	for i := range storage.tasks {
		taskDur := time.Duration(i+1) * time.Second // for diff task run freq
		storage.tasks[i] = task.NewTask(i, taskDur)
	}

	return storage, time.Duration(len(storage.tasks)+1) * time.Second
}

type Storage struct {
	lock  sync.Mutex
	tasks []*task.Task
}

var ErrCannotAcquire = errors.New("cannot acquire tasks storage")

func (t *Storage) Tasks() ([]*task.Task, error) {
	if !t.lock.TryLock() {
		return nil, ErrCannotAcquire
	}

	time.Sleep(100 * time.Millisecond) // work simulation

	tasks := make([]*task.Task, len(t.tasks))
	for i := range tasks {
		task := &t.tasks[i]
		tasks[i] = *task
	}

	t.lock.Unlock()

	return tasks, nil
}
