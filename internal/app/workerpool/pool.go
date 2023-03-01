package workerpool

import (
	"log"
	"sync"
	"time"

	"github.com/ShishkovEM/amazing-shortener/internal/app/interfaces"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
)

type DeletionPool struct {
	Tasks   []*models.DeletionTask
	Workers []*DeletionWorker

	concurrency   int
	collector     chan *models.DeletionTask
	runBackground chan bool
	target        interfaces.ProcessorTarget
	wg            sync.WaitGroup
}

func NewDeletionPool(tasks []*models.DeletionTask, concurrency int, target interfaces.ProcessorTarget) *DeletionPool {
	return &DeletionPool{
		Tasks:       tasks,
		concurrency: concurrency,
		collector:   make(chan *models.DeletionTask, 1000),
		target:      target,
	}
}

func (dp *DeletionPool) AddTask(task *models.DeletionTask) {
	dp.collector <- task
}

func (dp *DeletionPool) RunBackground() {
	go func() {
		for {
			log.Print("⌛ Waiting for tasks to come in ...\n")
			time.Sleep(10 * time.Second)
		}
	}()

	for i := 1; i <= dp.concurrency; i++ {
		worker := NewDeletionWorker(dp.collector, i, dp.target)
		dp.Workers = append(dp.Workers, worker)
		go worker.StartBackground()
	}

	for i := range dp.Tasks {
		dp.collector <- dp.Tasks[i]
	}

	dp.runBackground = make(chan bool)
	<-dp.runBackground
}

func (dp *DeletionPool) Stop() {
	for i := range dp.Workers {
		dp.Workers[i].Stop()
	}
	dp.runBackground <- true
}
