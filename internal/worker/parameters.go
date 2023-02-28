package worker

import (
	"sync/atomic"
	"time"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

var workerQty = 5
var maxPeriodToStorageSec = 3

const retryWorkerConnectExternalSeconds = 120

// multiplyes on workerQty
var newTaskQtyMultiplicator = 15
var limitTaskQtyToStorageMultiplicator = 10
var incomeChanMult = 20
var resultChanMult = 20

// max tasks qty limit
var maxTasks = workerQty * newTaskQtyMultiplicator

type ListOfTasks struct {
	ClientURL     string
	chTask        chan storage.Order
	chDelay       chan int
	chWorkers     chan chan time.Time
	arrayChan     []chan time.Time
	chResult      chan storage.Order
	workersPoused atomic.Bool  // is banned
	tasksReady    atomic.Int32 // tasks done and need to be uploaded new same qty
	taskTotal     int          // total tasks in memory (collector used only)
}
