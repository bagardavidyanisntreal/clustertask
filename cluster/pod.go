package cluster

import (
	"log"
	"sync"
	"time"

	"github.com/bagardavidyanisntreal/clustertask/task"
)

type (
	tasker interface {
		Tasks() ([]*task.Task, error)
	}
	podTasks interface {
		SetTasksByPodID(podID int, val []*task.Task)
		TasksByPodID(podID int) []*task.Task
	}
)

var cacheOnce sync.Once

func NewPod(id int, podsTotalCnt int, cacheTasksDur time.Duration, tasker tasker, podTasks podTasks) *Pod {
	pod := &Pod{
		cacheTicker:  time.NewTicker(cacheTasksDur),
		podsTotalCnt: podsTotalCnt,

		id:       id,
		tasker:   tasker,
		podTasks: podTasks,
	}

	cacheOnce.Do(func() {
		if err := pod.cacheTasks(); err != nil {
			log.Printf("start pod %d cache tasks failure: %q", pod.id, err)
		}
	})

	return pod
}

type Pod struct {
	cacheTicker  *time.Ticker
	podsTotalCnt int

	id       int
	tasker   tasker
	podTasks podTasks
}
