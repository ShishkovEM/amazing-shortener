package service

import (
	"encoding/json"
	"net/http"

	"github.com/ShishkovEM/amazing-shortener/internal/app/responses"

	"github.com/go-chi/chi/v5"
)

func (ls *LinkService) UserLinkRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/urls", ls.getLinksByUserIDHandler) // Получение ссылок, созданных полизователем
	return r
}

func (ls *LinkService) getLinksByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	rawUserID := r.Context().Value("userID")
	var userID uint64

	switch uidType := rawUserID.(type) {
	case uint64:
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
