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
	podsTotal     = 4
	taskTotal     = 100
	cacheTasksDur = 20 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	taskStorage := taskstorage.NewStorage(taskTotal)
	tasksByPod := imdbstorage.NewIMDB(podsTotal)

	pods := make([]*cluster.Pod, podsTotal)
	for i := range pods {
		pods[i] = cluster.NewPod(
			i,
			podsTotal,
			cacheTasksDur,
			taskStorage,
			tasksByPod,
		)
	}

	done := make(chan struct{}, len(pods)*2)

	for _, pod := range pods {
		pod := pod
		go func() {
			pod.CacheTasks(ctx)
			done <- struct{}{}
		}()

		go func() {
			pod.RunTasks(ctx)
			done <- struct{}{}
		}()
	}

	for i := 0; i < cap(done); i++ {
		<-done
	}
}
