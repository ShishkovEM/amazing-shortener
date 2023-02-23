package workerpool

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"

	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
)

type DeletionWorker struct {
	ID         int
	taskChan   chan *DeletionTask
	quit       chan bool
	connection *pgx.Conn
}

func NewDeletionWorker(channel chan *DeletionTask, ID int, DB *models.DB) *DeletionWorker {
	conn, err := DB.GetConn(context.Background())
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

func (dw *DeletionWorker) processDeletion(dt *DeletionTask) {
	log.Printf("Worker %d processes deletion of url %s\n", dw.ID, dt.urlToDelete)

	_, err := dw.connection.Exec(context.Background(), "UPDATE urls SET is_deleted = true WHERE short_uri = $1 AND user_id = $2", dt.urlToDelete, dt.userID)
	if err != nil {
		log.Printf("error updating URL: %v", err)
	}
}
