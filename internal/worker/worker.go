package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RunWorker operates handling external cervice.
// Get new task from task channel, returns completed ones in channel, or dropped task if status not changed.
func RunWorker(list *ListOfTasks, wg *sync.WaitGroup) {
	defer wg.Done()

	stopchan := make(chan time.Time, 1)
	list.chWorkers <- stopchan

	for {
		select {
		// stopchan should be chacked priority
		case deadline, ok := <-stopchan:
			if !ok {
				log.Println("worker: worker stopped")
				return
			}
			// check is deadline still actual
			nowtime := time.Now()
			if nowtime.Before(deadline) {
				log.Println("worker: worker paused for ", deadline.Sub(nowtime))
				// wait
				<-time.After(deadline.Sub(nowtime))
			}
		default:
			select {
			// in case of no tasks in channel we need keep listening stopchan
			case deadline, ok := <-stopchan:
				if !ok {
					log.Println("worker: worker stopped")
					return
				}
				// check is deadline still actual
				nowtime := time.Now()
				if nowtime.Before(deadline) {
					log.Println("worker: worker paused for ", deadline.Sub(nowtime))
					// wait
					<-time.After(deadline.Sub(nowtime))
				}
			case task := <-list.chTask:
				// log.Println("worker: next work:->", task.Number)
				adr := fmt.Sprintf("%s/api/orders/%s", list.ClientURL, task.Number)

				request, err := http.NewRequest(http.MethodGet, adr, nil)
				if err != nil {
					log.Println("worker: request build err:", err)
				}

				client := &http.Client{
					Timeout: 10 * time.Second,
				}

				resp, err := client.Do(request)
				if err != nil {
					log.Printf("worker: client err: %v. Retry after %d seconds\n", err, retryWorkerConnectExternalSeconds)
					// wait if external service doesn't respond
					<-time.After(time.Duration(retryWorkerConnectExternalSeconds) * time.Second)
					break
				}
				defer resp.Body.Close()

				// handle answers
				switch resp.StatusCode {
				case 200:
					if resp.Header.Get("Content-Type") != "application/json" {
						log.Println("worker: wrong header JSON in status 200")
					}

					oldStatus := task.Status

					// JSON
					result, _ := io.ReadAll(resp.Body)
					if err := json.Unmarshal(result, &task); err != nil {
						log.Println("worker: can't unmarshal JSON in answer")
					}

					// +1 jobs compleate counter
					list.tasksReady.Add(1)

					log.Println("GetEXTRA = ", task)

					// skip task w/o changes in status
					if oldStatus != task.Status {
						list.chResult <- task
					}

				case 429:

					delaystring := resp.Header.Get("Retry-After")
					// log.Printf("worker: 429: BAN, detay = %s", delaystring)
					delay, er := strconv.Atoi(delaystring) // in seconds?
					if er != nil {
						log.Println("worker: can't count seconds: got 429 with delay = ", delaystring)

						// wait 1 sec if status code is 429 but cant read time to wait
						delay = 1
					}

					// alert queue
					list.chDelay <- delay

					// task back to queue
					list.chTask <- task

				case 500:

					log.Println("worker: 500: got code 500 from service. wait 5 sec and go again for some reason")

					// in case of temporary error
					list.chDelay <- 5

					// task back to queue
					list.chTask <- task
				default:

					log.Println("worker: some problem with extra service. status code: ", resp.StatusCode)

					// task back to queue
					list.chTask <- task
				}
				log.Printf("worker: work result:->%s STATUS:%s", task.Number, task.Status)
			}

		}
	}
}
