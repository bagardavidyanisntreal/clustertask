package cluster

import (
	"context"
)

func (p *Pod) RunTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-p.ready:
			if !ok {
				return
			}

			tasks := p.podTasks.TasksByPodID(p.id)
			for _, task := range tasks {
				task := task
				go func() {
					task.Run(p.id, ctx)
				}()
			}
		}
	}
}
