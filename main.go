package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bagardavidyanisntreal/clustertask/cluster"
	imdbstorage "github.com/bagardavidyanisntreal/clustertask/storage/imdb"
	taskstorage "github.com/bagardavidyanisntreal/clustertask/storage/task"
)

const (
	podsTotal = 5
	taskTotal = 10
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	taskStorage, maxDur := taskstorage.NewStorage(taskTotal)
	tasksByPod := imdbstorage.NewIMDB(podsTotal)

	maxDur = maxDur + time.Second // for sure

	pods := make([]*cluster.Pod, podsTotal)
	for i := range pods {
		pods[i] = cluster.NewPod(
			ctx,
			i,
			podsTotal,
			maxDur,
			taskStorage,
			tasksByPod,
		)
	}

	var wg sync.WaitGroup
	wg.Add(len(pods))

	for _, pod := range pods {
		pod := pod
		go func() {
			pod.RunTasks(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
}
