package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// StartWorkers run worker and collector(1) goroutines, collects tasks and results in queue
func StartWorkers(ctx context.Context, con *storage.ServerConfigs) {
	// ctxworkers, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup

	TaskList := ListOfTasks{
		ClientURL: con.ExternalService,
		chTask:    make(chan storage.Order, incomeChanMult*workerQty),
		chDelay:   make(chan int),
		chWorkers: make(chan chan time.Time, 1),
		arrayChan: make([]chan time.Time, 0),
		chResult:  make(chan storage.Order, resultChanMult*workerQty),
	}

	// log.Println("adr con:", con.ExternalService, " takelistClURL:", TaskList.ClientURL)

	// start workers
	for i := 0; i < workerQty; i++ {
		wg.Add(1)
		go RunWorker(&TaskList, &wg)
	}

	// collect handle channels from every worker
	for i := 0; i < workerQty; i++ {
		select {
		case newch := <-TaskList.chWorkers:
			TaskList.arrayChan = append(TaskList.arrayChan, newch)
		case <-time.After(3 * time.Second):
			log.Println("can't recieve chan from worker")
		}
	}
	log.Printf("%d workers starts", len(TaskList.arrayChan))

	// start collector
	go TaskList.Collector(ctx)
	log.Println("collector starts")

	// listener
	for {
		select {
		case <-ctx.Done():

			// send stop to workers
			for _, ch := range TaskList.arrayChan {
				close(ch)
			}
			// wait workers to terminate
			wg.Wait()
			// signal to collector to save and terminate
			close(TaskList.chResult)
			return
		case delay := <-TaskList.chDelay:
			if delay == 0 {
				continue
			}

			// log.Println("BAN detected")

			// if BAN activated already
			da := TaskList.workersPoused.CompareAndSwap(false, true)
			if da {

				// set deadline
				deadline := time.Now().Add(time.Duration(delay) * time.Second)
				for _, ch := range TaskList.arrayChan {
					ch <- deadline
				}

				// clear worker chans after delay in case of worker lags
				go TaskList.clearWorckerChan(ctx, delay)
			}
		}
	}
}

// clearWorckerChan change status operates or wait T/O and clear channels
func (list *ListOfTasks) clearWorckerChan(ctx context.Context, delay int) {
	// log.Println("run waiter, delay = ", delay)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(delay) * time.Second):

			// release locker marker
			list.workersPoused.Store(false)

			for i := 0; i < len(list.arrayChan); i++ {
				select {
				// empty workers channels
				case <-list.arrayChan[i]:
				case <-time.After(10 * time.Microsecond):
				}
			}
		}
	}
}
