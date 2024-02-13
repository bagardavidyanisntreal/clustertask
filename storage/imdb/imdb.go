package imdb

import (
	"sync"

	"github.com/bagardavidyanisntreal/clustertask/dto"
)

func NewIMDB(ballast int) *IMDB {
	return &IMDB{
		tasksByPodID: make(map[int][]*dto.Task, ballast),
	}
}

type IMDB struct {
	lock         sync.Mutex
	tasksByPodID map[int][]*dto.Task
}

func (i *IMDB) SetTasksByPodID(podID int, val []*dto.Task) {
	i.lock.Lock()
	i.tasksByPodID[podID] = val
	i.lock.Unlock()
}
func (i *IMDB) TasksByPodID(podID int) []*dto.Task {
	i.lock.Lock()
	defer i.lock.Unlock()

	tasks, ok := i.tasksByPodID[podID]
	if !ok {
		return nil
	}

	return tasks
}
