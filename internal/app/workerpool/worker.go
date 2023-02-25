package workerpool

import (
	"context"
	"log"

	"github.com/ShishkovEM/amazing-shortener/internal/app/interfaces"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"

	"github.com/jackc/pgx/v4"
)

type DeletionWorker struct {
	ID               int
	taskChan         chan *models.DeletionTask
	quit             chan bool
	activeConnection *pgx.Conn
	querier          interfaces.Queriable
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
	err := dw.activeConnection.Close(context.Background())
	if err != nil {
		log.Printf("error closing DB querier: %v", err)
	}
	go func() {
		dw.quit <- true
	}()
}

func (dw *DeletionWorker) processDeletion(dt *models.DeletionTask) {
	log.Printf("Worker %d processes deletion of url %s\n", dw.ID, dt.URLToDelete)

	if dw.activeConnection == nil {
		var conn *pgx.Conn
		conn, err := dw.querier.GetConn(context.Background())
		if err != nil {
			log.Fatalf("error establising connection: %v", err)
		}
		dw.activeConnection = conn
	}

	_, err := dw.activeConnection.Exec(context.Background(), "UPDATE urls SET is_deleted = true WHERE short_uri = $1 AND user_id = $2", dt.URLToDelete, dt.UserID)
	if err != nil {
		log.Printf("error updating URL: %v", err)
	}
}
