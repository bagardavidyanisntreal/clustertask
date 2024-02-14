package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/bagardavidyanisntreal/clustertask/cluster"
	imdbstorage "github.com/bagardavidyanisntreal/clustertask/storage/imdb"
	taskstorage "github.com/bagardavidyanisntreal/clustertask/storage/task"
)

const (
	podsTotal = 5
	taskTotal = 15
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	taskStorage, maxDur := taskstorage.NewStorage(taskTotal)
	tasksByPod := imdbstorage.NewIMDB(podsTotal)
	defer tasksByPod.Close()

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

	done := make(chan struct{}, len(pods))

	for _, pod := range pods {
		pod := pod
		go func() {
			pod.RunTasks(ctx)
			done <- struct{}{}
		}()
	}

	for range pods {
		<-done
	}
}
