package workerpool

import (
	"context"
	"log"

	"github.com/ShishkovEM/amazing-shortener/internal/app/interfaces"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"

	"github.com/jackc/pgx/v4"
)

type DeletionWorker struct {
	ID         int
	taskChan   chan *models.DeletionTask
	quit       chan bool
	connection *pgx.Conn
}

func NewDeletionWorker(channel chan *models.DeletionTask, ID int, target interfaces.ProcessorTarget) *DeletionWorker {
	conn, err := target.GetConn(context.Background())
	if err != nil {
		log.Fatal("Error establishing connection to database")
	}
	return &DeletionWorker{
		ID:         ID,
		taskChan:   channel,
		quit:       make(chan bool),
		connection: conn,
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
	err := dw.connection.Close(context.Background())
	if err != nil {
		log.Printf("error closing DB connection: %v", err)
	}
	go func() {
		dw.quit <- true
	}()
}

func (dw *DeletionWorker) processDeletion(dt *models.DeletionTask) {
	log.Printf("Worker %d processes deletion of url %s\n", dw.ID, dt.URLToDelete)

	_, err := dw.connection.Exec(context.Background(), "UPDATE urls SET is_deleted = true WHERE short_uri = $1 AND user_id = $2", dt.URLToDelete, dt.UserID)
	if err != nil {
		log.Printf("error updating URL: %v", err)
	}
}
