package worker

import (
	"sync/atomic"
	"time"

	"github.com/kormiltsev/gophermartbonus/internal/storage"
)

// workerQty is max worker operates same time.
var workerQty = 5

// maxPeriodToStorageSec is maximum time duration between upload tasks completed and get new ones from Postgress.
var maxPeriodToStorageSec = 3

// retryWorkerConnectExternalSeconds sets time limit to reconnect external service in case of connection error.
const retryWorkerConnectExternalSeconds = 120

// Vars used by queue, multiplyes on workerQty.
var (
	newTaskQtyMultiplicator            = 15
	limitTaskQtyToStorageMultiplicator = 10
	incomeChanMult                     = 20
	resultChanMult                     = 20
)

// Max tasks qty limit.
var maxTasks = workerQty * newTaskQtyMultiplicator

// ListOfTasks manages queue of tasks.
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
