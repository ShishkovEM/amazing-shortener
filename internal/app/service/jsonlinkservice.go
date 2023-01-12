package service

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (ls *LinkService) RestRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/shorten", ls.createLinkJSONHandler) // Создание новой сокращённой ссылки
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
	// Структуры для запроса и ответа
	type RequestLink struct {
		URL string `json:"url" validate:"required,url"`
	}

	type ResponseShortLink struct {
		Result string `json:"result"`
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

	var rl RequestLink
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

	short := ls.store.CreateLink(rl.URL)
	log.Printf("created short id: %s\n", short)
	renderJSON(w, ResponseShortLink{Result: ls.baseURL + short})
}
