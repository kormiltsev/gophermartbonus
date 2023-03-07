package worker

import (
	"context"
	"log"
	"time"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// Collector collect results from queue. Push to DB results and get new tasks
func (list *ListOfTasks) Collector(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(maxPeriodToStorageSec) * time.Second)
	// log.Println("run collector")

	// request PG for first tasks
	newtasks, err := storage.UpdateRows(ctx, nil, newTaskQtyMultiplicator*workerQty)
	if err == nil {
		for _, v := range newtasks {
			list.chTask <- v
		}
	}

	answers := make([]storage.Order, 0)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Println("collector dead")
			return
		case res, ok := <-list.chResult:
			if !ok {
				_, _ = storage.UpdateRows(ctx, answers, 0)
				ticker.Stop()
				log.Println("collector dead")
				return
			}
			// log.Println("collector accepts: ", res.Number)
			answers = append(answers, res)

			if int(list.tasksReady.Load()) > limitTaskQtyToStorageMultiplicator*workerQty {
				tasksDone := int(list.tasksReady.Swap(0))
				log.Printf("start push postgres by limit, tasks done = %d, len(answers)=%d", tasksDone, len(answers))

				// if error from DB so no new task. work continue on old data
				newtasks, err := storage.UpdateRows(ctx, answers, tasksDone)
				if err == nil {
					for _, v := range newtasks {
						list.chTask <- v
					}
				}

				// log.Println("finish postgres by limit")
				if len(answers) != 0 {
					answers = answers[:0]
				}
			}
		case <-ticker.C:

			tasksDone := int(list.tasksReady.Swap(0))
			list.taskTotal -= tasksDone

			// skip if no word done
			if tasksDone == 0 && maxTasks-list.taskTotal == 0 {
				return
			}

			if len(answers) != 0 {
				log.Println("start postgres by timer, tasks ready =", len(answers))
			}

			// run request to PG
			// if error from DB so no new task. work continue on old data
			newtasks, err := storage.UpdateRows(ctx, answers, maxTasks-list.taskTotal)
			if err == nil {
				for _, v := range newtasks {
					list.chTask <- v
				}
			}

			list.taskTotal += len(newtasks)

			// log.Println("finish postgres by timer")
			if len(answers) != 0 {
				answers = answers[:0]
			}

		}
	}
}
