package service

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShishkovEM/amazing-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

func TestLinkService_CreateLinkJSONHandlerPositive(t *testing.T) {
	linkStorage := storage.NewLinkStore("")
	ls := NewLinkService(linkStorage, "http://"+linkServiceHost+":"+linkServicePort+"/")
	r := chi.NewRouter()
	r.Mount("/api", ls.RestRoutes())
	reqBody := bytes.NewBuffer([]byte("{\"url\":\"" + testedLongURL + "\"}"))

	req, err := http.NewRequest("POST", "/api/shorten", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusCreated
		mediaTypeOK := w.Header().Get("Content-Type") == "application/json"
		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Contains(string(p), linkServiceHost+":"+linkServicePort)

		return statusOK && mediaTypeOK && pageOK
	})
}

func TestLinkService_CreateLinkJSONHandlerInvalidURL(t *testing.T) {
	linkStorage := storage.NewLinkStore("")
	ls := NewLinkService(linkStorage, "http://"+linkServiceHost+":"+linkServicePort+"/")
	r := chi.NewRouter()
	r.Mount("/api", ls.RestRoutes())
	reqBody := bytes.NewBuffer([]byte("{\"url\":\"" + testedInvalidURL + "\"}"))

	req, err := http.NewRequest("POST", "/api/shorten", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusBadRequest
		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Contains(string(p), "Invalid URL creation request handled")

		return statusOK && pageOK
	})
}
