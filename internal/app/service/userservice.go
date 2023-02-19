package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/middleware"
	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"

	"github.com/go-chi/chi/v5"
)

func (ls *LinkService) UserLinkRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/urls", ls.getLinksByUserIDHandler)  // Получение ссылок, созданных полизователем
	r.Delete("/urls", ls.deleteUserURLsHandler) // Удаление ссылок
	return r
}

func (ls *LinkService) getLinksByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	rawUserID := r.Context().Value(middleware.ContextKeyUserID)
	var userID uint32

	switch uidType := rawUserID.(type) {
	case uint32:
		userID = uidType
	}

	userLinks := ls.store.GetLinksByUserID(userID)

	var response []responses.ResponseShortOriginalLink

	for _, link := range userLinks {
		response = append(response, responses.ResponseShortOriginalLink{ShortURL: ls.baseURL + link.Short, OriginalURL: link.Original})
	}

	responseBytes, _ := json.Marshal(response)

	if len(userLinks) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, err := w.Write(responseBytes)
	if err != nil {
		return
	}
}

func (ls *LinkService) deleteUserURLsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling link delete via #deleteUserURLsHandler at %s\n", req.URL.Path)
	w.Header().Set("content-type", "application/json")

	rawUserID := req.Context().Value(middleware.ContextKeyUserID)
	var userID uint32

	switch uidType := rawUserID.(type) {
	case uint32:
		userID = uidType
	}

	var request []string

	b, _ := io.ReadAll(req.Body)
	err := json.Unmarshal(b, &request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, `{"error":"Invalid body"}`, http.StatusBadRequest)
		return
	}

	if len(request) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, `{"error":"Body must be not empty array"}`, http.StatusBadRequest)
		return
	}

	err = ls.store.DeleteUserRecordsByShortURLs(userID, request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, `{"error":"Failed to delete user records by short url"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
