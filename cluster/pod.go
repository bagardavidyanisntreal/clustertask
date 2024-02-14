package cluster

import (
	"context"
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
		Ready() chan struct{}
	}
)

func NewPod(
	ctx context.Context,
	id int,
	podsTotalCnt int,
	cacheDur time.Duration,
	tasker tasker,
	podTasks podTasks,
) *Pod {
	pod := &Pod{
		ticker:       time.NewTicker(cacheDur),
		podsTotalCnt: podsTotalCnt,

		id:       id,
		tasker:   tasker,
		podTasks: podTasks,
	}

	go pod.runCaching(ctx)

	return pod
}

type Pod struct {
	ticker       *time.Ticker
	podsTotalCnt int

	id       int
	tasker   tasker
	podTasks podTasks
}
