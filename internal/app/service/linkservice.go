package service

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ShishkovEM/amazing-shortener/internal/app/middleware"
	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

type LinkService struct {
	baseURL string
	store   *storage.LinkStore
}

func NewLinkService(store *storage.LinkStore, baseURL string) *LinkService {
	ls := &LinkService{
		baseURL: baseURL,
		store:   store,
	}
	return ls
}

func (ls *LinkService) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ls.createLinkHandler) // Создание новой сокращённой ссылки

	r.Get("/{short}", ls.getLinkHandler) // Восстановление ссылки

	return r
}

func (ls *LinkService) createLinkHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling link create at %s\n", req.URL.Path)

	rawUserID := req.Context().Value(middleware.ContextKeyUserID)
	var userID uint64

	switch uidType := rawUserID.(type) {
	case uint64:
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

	short, err := ls.store.CreateLink(string(LongURL), userID)
	if err != nil {
		log.Printf("Error creating link: %s\n", err)
		return
	}

	responseBody, err := ls.store.GetLink(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("created short id: %s\n", responseBody.Short)

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(ls.baseURL + responseBody.Short))
	if err != nil {
		log.Printf("Error writing response body at createLinkHandler: %s\n", err)
	}
}

func (ls *LinkService) getLinkHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get link at %s\n", req.URL.Path)

	link, err := ls.store.GetLink(strings.Trim(req.URL.Path, "/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("expanded original link: %s\n", link.Original)

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.Header().Set("Location", link.Original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func isValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}
	u, err := url.Parse(input)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}
