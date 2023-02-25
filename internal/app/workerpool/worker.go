package workerpool

import (
	"context"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"log"

	"github.com/ShishkovEM/amazing-shortener/internal/app/interfaces"
)

type DeletionWorker struct {
	ID       int
	taskChan chan *models.DeletionTask
	quit     chan bool
	querier  interfaces.Queriable
}

func NewDeletionWorker(channel chan *models.DeletionTask, ID int, querier interfaces.Queriable) *DeletionWorker {
	return &DeletionWorker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
		querier:  querier,
	}
}

func (dw *DeletionWorker) StartBackground() {
	log.Printf("Starting worker %d\n", dw.ID)

	for {
		select {
		case deletionTask := <-dw.taskChan:
			dw.processDeletion(deletionTask)
		case <-dw.quit:
			return
		}
	}
}

func (dw *DeletionWorker) Stop() {
	log.Printf("Closing worker %d\n", dw.ID)
	err := dw.querier.Close()
	if err != nil {
		log.Printf("error closing DB querier: %v", err)
	}
	go func() {
		dw.quit <- true
	}()
}

func (dw *DeletionWorker) processDeletion(dt *models.DeletionTask) {
	log.Printf("Worker %d processes deletion of url %s\n", dw.ID, dt.UrlToDelete)

	q, err := dw.querier.GetQuerier()
	if err != nil {
		log.Printf("error getting execer: %v", err)
	}

	_, err = q.Exec(context.Background(), "UPDATE urls SET is_deleted = true WHERE short_uri = $1 AND user_id = $2", dt.UrlToDelete, dt.UserID)
	if err != nil {
		log.Printf("error updating URL: %v", err)
	}
}
