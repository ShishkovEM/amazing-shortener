package service

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/middleware"
	"github.com/ShishkovEM/amazing-shortener/internal/app/requests"
	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"

	"github.com/go-chi/chi/v5"
)

func (ls *LinkService) RestRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/shorten", ls.createLinkJSONHandler)         // Создание новой сокращённой ссылки
	r.Post("/shorten/batch", ls.createLinksBatchHandler) // Пакетное создание ссылок
	r.Delete("/user/urls", ls.deleteUserURLsHandler)     // Удаление ссылок
	return r
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(js)
	if err != nil {
		log.Printf("Error rendering JSON: %v\n", err)
		return
	}
}

func (ls *LinkService) createLinkJSONHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling link create via #createLinkJSONHandler at %s\n", req.URL.Path)

	rawUserID := req.Context().Value(middleware.ContextKeyUserID)
	var userID uint32

	switch uidType := rawUserID.(type) {
	case uint32:
		userID = uidType
	}

	// Проверяем, что на вход получен JSON
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("Error parsing mediatype: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		log.Printf("Error parsing mediatype. Expected \"application/json\" Got: %s\n", mediaType)
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var rl requests.RequestLink
	if err := dec.Decode(&rl); err != nil {
		log.Printf("Error decoding link request: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("recieved long url: %s\n", rl.URL)

	if !isValidURL(rl.URL) {
		log.Printf("Handled invlid URL: %s\n", rl.URL)
		http.Error(w, "Invalid URL creation request handled. Input URL: "+rl.URL, http.StatusBadRequest)
		return
	}

	short, err := ls.store.CreateLink(rl.URL, userID)
	if err != nil {
		log.Printf("Error creating link: %s\n", err)
		return
	}
	log.Printf("created short id: %s\n", short)
	renderJSON(w, responses.ResponseShortLink{Result: ls.baseURL + short})
}

func (ls *LinkService) createLinksBatchHandler(w http.ResponseWriter, req *http.Request) {
	var request []requests.RequestLinksBatch
	var result []responses.ResponseLinksBatch

	rawUserID := req.Context().Value(middleware.ContextKeyUserID)
	var userID uint32

	switch uidType := rawUserID.(type) {
	case uint32:
		userID = uidType
	}

	b, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(b, &request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(request) == 0 {
		http.Error(w, `{"error":"URLs in body are required"}`, http.StatusBadRequest)
		return
	}

	for _, data := range request {
		shortURL, err := ls.store.CreateLink(data.OriginalURL, userID)
		if err != nil {
			log.Printf("Error creating link in batch")
			return
		}
		result = append(result, responses.ResponseLinksBatch{CorrelationID: data.CorrelationID, ShortURL: ls.baseURL + shortURL})
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseBytes, _ := json.Marshal(result)

	_, err = w.Write(responseBytes)
	if err != nil {
		return
	}
}
