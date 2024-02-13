package cluster

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bagardavidyanisntreal/clustertask/dto"
)

type (
	tasker interface {
		Tasks() ([]*dto.Task, error)
	}
	podTasks interface {
		SetTasksByPodID(podID int, val []*dto.Task)
		TasksByPodID(podID int) []*dto.Task
	}
)

var cacheOnce sync.Once

func NewPod(id int, podsTotalCnt int, cacheTasksDur time.Duration, tasker tasker, podTasks podTasks) *Pod {
	pod := &Pod{
		cacheTasksDur: cacheTasksDur,
		podsTotalCnt:  podsTotalCnt,

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
	cacheTasksDur time.Duration
	podsTotalCnt  int

	id       int
	tasker   tasker
	podTasks podTasks
}

func (p *Pod) ID() int {
	return p.id
}

func (p *Pod) CacheTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("stopping pod %d\n", p.id)
			return
		case <-time.After(p.cacheTasksDur):
			if err := p.cacheTasks(); err != nil {
				log.Printf("pod %d cache tasks failure: %q", p.id, err)
			}
		}
	}
}

func (p *Pod) cacheTasks() error {
	tasks, err := p.tasker.Tasks()
	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %w", err)
	}

	if p.podsTotalCnt > len(tasks) {
		return p.oneTaskPerPod(tasks)
	}

	batch := len(tasks) / p.podsTotalCnt
	var podID int

	for {
		if len(tasks) < batch || batch == 0 {
			batch = len(tasks)
		}

		process := tasks[:batch]
		if len(process) == 0 {
			break
		}

		p.podTasks.SetTasksByPodID(podID, process)
		log.Printf("succesfully cached %d tasks for pod %d", len(process), podID)

		tasks = tasks[batch:]
		podID++
	}

	return nil
}

func (p *Pod) oneTaskPerPod(tasks []*dto.Task) error {
	for podID := 0; podID < p.podsTotalCnt; podID++ {
		if podID > len(tasks)-1 {
			break
		}

		singleTask := tasks[podID]
		p.podTasks.SetTasksByPodID(podID, []*dto.Task{singleTask})
		log.Printf("succesfully cached task %d for pod %d", podID, podID)
	}

	return nil
}

func (p *Pod) RunTasks(ctx context.Context) {
	tasks := p.podTasks.TasksByPodID(p.id)
	done := make(chan struct{}, len(tasks))

	for _, task := range tasks {
		task := task
		go func() {
			task.Run(p.id, ctx)
			done <- struct{}{}
		}()
	}

	for i := 0; i < cap(done); i++ {
		<-done
	}
}
