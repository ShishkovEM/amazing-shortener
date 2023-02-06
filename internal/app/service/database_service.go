package service

import (
	"net/http"
	"sync"

	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

type DBService struct {
	sync.Mutex

	Db *storage.DB
}

func (dbs *DBService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", dbs.ping())
	return r
}

func NewDataBaseService(store *storage.DB) *DBService {
	dbs := &DBService{
		Db: store,
	}
	return dbs
}

func (dbs *DBService) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer dbs.Db.Close()

		conn, err := dbs.Db.GetConn(r.Context())
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
