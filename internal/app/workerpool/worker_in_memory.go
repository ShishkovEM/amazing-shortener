package workerpool

import (
	"log"

	"github.com/ShishkovEM/amazing-shortener/internal/app/interfaces"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
)

type MemDeletionWorker struct {
	ID       int
	taskChan chan *models.DeletionTask
	quit     chan bool
	target   interfaces.InMemoryLinkStorage
}

func NewMemDeletionWorker(channel chan *models.DeletionTask, ID int, target interfaces.InMemoryLinkStorage) *MemDeletionWorker {
	return &MemDeletionWorker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
		target:   target,
	}
}

func (mdw *MemDeletionWorker) StartBackground() {
	log.Printf("Starting worker %d\n", mdw.ID)

	for {
		select {
		case deletionTask := <-mdw.taskChan:
			mdw.processDeletion(deletionTask)
		case <-mdw.quit:
			return
		}
	}
}

func (mdw *MemDeletionWorker) Stop() {
	log.Printf("Closing worker %d\n", mdw.ID)

	go func() {
		mdw.quit <- true
	}()
}

func (mdw *MemDeletionWorker) processDeletion(dt *models.DeletionTask) {
	log.Printf("Worker %d processes deletion of url %s\n", mdw.ID, dt.URLToDelete)

	mdw.target.DeleteOne(dt)
}
