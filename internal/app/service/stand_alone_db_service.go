package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ShishkovEM/amazing-shortener/internal/app/middleware"
	"github.com/ShishkovEM/amazing-shortener/internal/app/models"
	"github.com/ShishkovEM/amazing-shortener/internal/app/repository"
	"github.com/ShishkovEM/amazing-shortener/internal/app/requests"
	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"

	"github.com/go-chi/chi/v5"
	"github.com/speps/go-hashids"
)

type StandAloneDBService struct {
	baseURL string
	store   *repository.DBLinkStorage
}

func NewStandAloneDBService(store *repository.DBLinkStorage, baseURL string) *StandAloneDBService {
	sadbs := &StandAloneDBService{
		store:   store,
		baseURL: baseURL,
	}
	return sadbs
}

func (sadbs *StandAloneDBService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", sadbs.createLinkHandler)                // Создание новой сокращённой ссылки
	r.Get("/{short}", sadbs.getLinkHandler)             // Восстановление ссылки
	r.Post("/api/shorten", sadbs.createLinkJSONHandler) // Создание новой сокращённой ссылки
	r.Get("/api/urls", sadbs.getLinksByUserIDHandler)   // Получение ссылок, созданных полизователем
	r.Get("/ping", sadbs.ping())

	return r
}

func (sadbs *StandAloneDBService) createLinkHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling link create at %s\n", req.URL.Path)

	rawUserID := req.Context().Value(middleware.ContextKeyUserID)
	var userID uint32

	switch uidType := rawUserID.(type) {
	case uint32:
		userID = uidType
	}

	LongURL, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("recieved long url: %s\n", LongURL)

	if !isValidURL(string(LongURL)) {
		http.Error(w, "Invalid URL creation request handled. Input URL: "+string(LongURL), http.StatusBadRequest)
		return
	}

	link := models.Link{
		Original: string(LongURL),
		Short:    shorten(),
		UserID:   userID,
	}

	sadbs.store.CreateLink(link.Short, link.Original, link.UserID)
	if err != nil {
		log.Printf("Error creating link: %s\n", err)
		return
	}

	log.Printf("created short id: %s\n", link.Short)

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(sadbs.baseURL + link.Short))
	if err != nil {
		log.Printf("Error writing response body at createLinkHandler: %s\n", err)
	}
}

func (sadbs *StandAloneDBService) getLinkHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get link at %s\n", req.URL.Path)

	link, err := sadbs.store.GetLink(strings.Trim(req.URL.Path, "/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("expanded original link: %s\n", link)

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (sadbs *StandAloneDBService) createLinkJSONHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling link create via #createLinkJSONHandler at %s\n", req.URL.Path)

	rawUserID := req.Context().Value(middleware.ContextKeyUserID)
	var userID uint32

	switch uidType := rawUserID.(type) {
	case uint32:
		userID = uidType
	}

	b, _ := io.ReadAll(req.Body)

	var request requests.RequestLink

	err := json.Unmarshal(b, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.URL == "" {
		http.Error(w, `{"error":"URL in body is required"}`, http.StatusBadRequest)
		return
	}

	link := models.Link{
		Original: request.URL,
		Short:    shorten(),
		UserID:   userID,
	}

	sadbs.store.CreateLink(link.Short, link.Original, link.UserID)
	response := responses.ResponseShortLink{Result: link.Short}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseBytes, _ := json.Marshal(response)

	_, err = w.Write(responseBytes)
	if err != nil {
		return
	}
}

func (sadbs *StandAloneDBService) getLinksByUserIDHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	rawUserID := req.Context().Value(middleware.ContextKeyUserID)
	var userID uint32

	switch uidType := rawUserID.(type) {
	case uint32:
		userID = uidType
	}

	userURLs := sadbs.store.GetLinksByUserID(userID)
	responseBytes, _ := json.Marshal(userURLs)

	if len(userURLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, err := w.Write(responseBytes)
	if err != nil {
		return
	}
}

func (sadbs *StandAloneDBService) ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer sadbs.store.Db.Close()

		conn, err := sadbs.store.Db.GetConn(r.Context())
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

func shorten() string {
	hd := hashids.NewData()
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	short, _ := h.Encode([]int{int(now.UnixMicro())})
	return short
}