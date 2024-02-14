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
		ready:        make(chan struct{}),

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
	ready        chan struct{}

	id       int
	tasker   tasker
	podTasks podTasks
}
