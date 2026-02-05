package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"lms_backend/pkg/broker"
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
						log.Printf("Redis error: %v", err)
						time.Sleep(time.Second)
						continue
					}
				}

				var task broker.Task
				if err := json.Unmarshal([]byte(res[1]), &task); err != nil {
					log.Printf("Unmarshal error: %v", err)
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
		log.Printf("Worker [%d] processing task: %s", id, t.Type)

		switch t.Type {
		case "EMAIL_CONFIRMATION":
			w.handleEmail(t.Payload)
		case "NEW_SUBMISSION":
			w.handleNewSubmission(t.Payload)
		}
	}
}

func (w *NotificationWorker) handleEmail(payload map[string]string) {
	log.Printf("Sending email to %s...", payload["email"])
	time.Sleep(time.Second * 1)
}

func (w *NotificationWorker) handleNewSubmission(payload map[string]string) {
	log.Printf("Notifying staff about new submission from student %s", payload["student_id"])
}
