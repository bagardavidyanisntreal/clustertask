package cluster

import (
	"context"
	"fmt"
	"log"

	"github.com/bagardavidyanisntreal/clustertask/task"
)

func (p *Pod) runCaching(ctx context.Context) {
	defer p.ticker.Stop()

	for {
		if err := p.cacheTasks(); err != nil {
			log.Printf("pod %d cache tasks failure: %q", p.id, err)
		}

		select {
		case <-ctx.Done():
			fmt.Printf("stopping pod %d\n", p.id)
			return
		case <-p.ticker.C:
			continue
		}
	}
}

const logSuccessCached = "[pod %d] successfully cached %d tasks for pod %d"

func (p *Pod) cacheTasks() error {
	tasks, err := p.tasker.Tasks()
	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %w", err)
	}

	if p.podsTotalCnt > len(tasks) {
		p.cacheOneTaskPerPod(tasks)
		return nil
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
		log.Printf(logSuccessCached, p.id, len(process), podID)

		tasks = tasks[batch:]
		podID++
	}

	return nil
}

func (p *Pod) cacheOneTaskPerPod(tasks []*task.Task) {
	for podID := 0; podID < p.podsTotalCnt; podID++ {
		if podID > len(tasks)-1 {
			break
		}

		singleTask := tasks[podID]
		batch := []*task.Task{singleTask}
		p.podTasks.SetTasksByPodID(podID, batch)

		log.Printf(logSuccessCached, p.id, len(batch), podID)
	}
}
