package workerpool

import (
	"log"
	"sync"
	"time"

	"github.com/ShishkovEM/amazing-shortener/internal/app/interfaces"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
)

type MemDeletionPool struct {
	Tasks   []*models.DeletionTask
	Workers []*MemDeletionWorker

	concurrency   int
	collector     chan *models.DeletionTask
	runBackground chan bool
	target        interfaces.InMemoryLinkStorage
	wg            sync.WaitGroup
}

func NewMemDeletionPool(tasks []*models.DeletionTask, concurrency int, target interfaces.InMemoryLinkStorage) *MemDeletionPool {
	return &MemDeletionPool{
		Tasks:       tasks,
		concurrency: concurrency,
		collector:   make(chan *models.DeletionTask, 1000),
		target:      target,
	}
}

func (mdp *MemDeletionPool) AddTask(task *models.DeletionTask) {
	mdp.collector <- task
}

func (mdp *MemDeletionPool) RunBackground() {
	go func() {
		for {
			log.Print("⌛ Waiting for tasks to come in ...\n")
			time.Sleep(10 * time.Second)
		}
	}()

	for i := 1; i <= mdp.concurrency; i++ {
		worker := NewMemDeletionWorker(mdp.collector, i, mdp.target)
		mdp.Workers = append(mdp.Workers, worker)
		go worker.StartBackground()
	}

	for i := range mdp.Tasks {
		mdp.collector <- mdp.Tasks[i]
	}

	mdp.runBackground = make(chan bool)
	<-mdp.runBackground
}

func (mdp *MemDeletionPool) Stop() {
	for i := range mdp.Workers {
		mdp.Workers[i].Stop()
	}
	mdp.runBackground <- true
}
