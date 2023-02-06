package service

import (
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type DBService struct {
	sync.Mutex

	DB *models.DB
}

func (dbs *DBService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", dbs.ping())
	return r
}

func NewDataBaseService(store *models.DB) *DBService {
	dbs := &DBService{
		DB: store,
	}
	return dbs
}

func (dbs *DBService) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer dbs.DB.Close()

		conn, err := dbs.DB.GetConn(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = conn.Ping(r.Context())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
