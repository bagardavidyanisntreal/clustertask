package imdb

import (
	"sync"

	"github.com/bagardavidyanisntreal/clustertask/task"
)

func NewIMDB(ballast int) *IMDB {
	return &IMDB{
		tasksByPodID: make(map[int][]*task.Task, ballast),
	}
}

type IMDB struct {
	lock         sync.Mutex
	tasksByPodID map[int][]*task.Task
}

func (i *IMDB) SetTasksByPodID(podID int, val []*task.Task) {
	i.lock.Lock()
	i.tasksByPodID[podID] = val
	i.lock.Unlock()
}
func (i *IMDB) TasksByPodID(podID int) []*task.Task {
	i.lock.Lock()
	defer i.lock.Unlock()

	tasks, ok := i.tasksByPodID[podID]
	if !ok {
		return nil
	}

	return tasks
}
