package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"lms_backend/pkg/broker"
	"lms_backend/pkg/logger"
)

type NotificationWorker struct {
	broker      *broker.EventBroker
	workerCount int
}

func NewNotificationWorker(b *broker.EventBroker, count int) *NotificationWorker {
	return &NotificationWorker{
		broker:      b,
		workerCount: count,
	}
}

func (w *NotificationWorker) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	taskChan := make(chan broker.Task, w.workerCount)

	for i := 0; i < w.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.startWorker(workerID, taskChan)
		}(i)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(taskChan)
				return
			default:
				res, err := w.broker.Subscribe(ctx, broker.NotificationEvent)
				if err != nil {
					select {
					case <-ctx.Done():
						return
					default:
						slog.Error("Redis error", logger.Err(err))
						time.Sleep(time.Second)
						continue
					}
				}

				var task broker.Task
				if err := json.Unmarshal([]byte(res[1]), &task); err != nil {
					slog.Error("Unmarshal error", logger.Err(err))
					continue
				}

				taskChan <- task
			}
		}
	}()

	wg.Wait()
	return nil
}

func (w *NotificationWorker) startWorker(id int, tasks <-chan broker.Task) {
	for t := range tasks {
		slog.Info("Worker processing task", slog.Int("worker_id", id), slog.String("task_type", t.Type))

		switch t.Type {
		case "EMAIL_CONFIRMATION":
			w.handleEmail(t.Payload)
		case "NEW_SUBMISSION":
			w.handleNewSubmission(t.Payload)
		}
	}
}

func (w *NotificationWorker) handleEmail(payload map[string]string) {
	slog.Info("Sending email", slog.String("email", payload["email"]))
	time.Sleep(time.Second * 1)
}

func (w *NotificationWorker) handleNewSubmission(payload map[string]string) {
	slog.Info("Notifying staff about new submission", slog.String("student_id", payload["student_id"]))
}
