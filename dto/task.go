package dto

import (
	"context"
	"fmt"
	"time"
)

func NewTask(id int, dur time.Duration) *Task {
	return &Task{
		id:   id,
		freq: dur,
	}
}

type Task struct {
	id   int
	freq time.Duration
}

func (t Task) Run(podID int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("task %d finish\n", t.id)
			return
		case <-time.After(t.freq):
			fmt.Printf("[pod %d]: running task %d\n", podID, t.id)
		}
	}
}
