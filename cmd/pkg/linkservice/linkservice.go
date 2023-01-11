package linkservice

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"

	"github.com/go-chi/chi/v5"
)

type LinkService struct {
	serviceAddress string
	store          *linkstore.LinkStore
}

func NewLinkService(store *linkstore.LinkStore, serviceAddress string) *LinkService {
	ls := &LinkService{
		serviceAddress: serviceAddress,
		store:          store,
	}
	return ls
}

func (ls *LinkService) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", ls.createLinkHandler) // Создание новой сокращённой ссылки

	r.Route("/{short}", func(r chi.Router) {
		r.Get("/", ls.getLinkHandler) // Восстановление ссылки
	})

	return r
}

func (ls *LinkService) createLinkHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling link create at %s\n", req.URL.Path)

	LongURL, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidURL(string(LongURL)) {
		http.Error(w, "Invalid URL creation request handled. Input URL: "+string(LongURL), http.StatusBadRequest)
		return
	}

	short := ls.store.CreateLink(string(LongURL))

	responseBody, err := ls.store.GetLink(short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(ls.serviceAddress + responseBody.Short))
	if err != nil {
		log.Print("Error writing response body at createLinkHandler")
	}
}

func (ls *LinkService) getLinkHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get link at %s\n", req.URL.Path)

	link, err := ls.store.GetLink(strings.Trim(req.URL.Path, "/"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

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
