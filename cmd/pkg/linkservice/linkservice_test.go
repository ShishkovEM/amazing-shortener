package linkservice

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"

	"github.com/go-chi/chi/v5"
)

const (
	testedLongURL    = "http://ya.ru/"
	testedInvalidURL = "not URL at all"
)

func testHTTPResponse(t *testing.T, r chi.Router, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func TestLinkServer_CreateLinkHandlerPositive(t *testing.T) {
	storage := linkstore.NewLinkStore()
	ls := NewLinkService(storage)
	r := ls.Routes()

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(testedLongURL)))
	if err != nil {
		t.Fatal(err)
	}

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusCreated
		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Contains(string(p), ls.Run(ls.serverAddress, ls.serverPort))

		return statusOK && pageOK
	})
}

func TestLinkServer_CreateLinkHandlerWithInvalidURL(t *testing.T) {
	storage := linkstore.NewLinkStore()
	ls := NewLinkService(storage)
	r := ls.Routes()

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(testedInvalidURL)))
	if err != nil {
		t.Fatal(err)
	}

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusBadRequest
		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Contains(string(p), "Invalid URL creation request handled. Input URL: ")

		return statusOK && pageOK
	})
}

func TestLinkServer_GetLinkHandlerPositive(t *testing.T) {
	storage := linkstore.NewLinkStore()
	ls := NewLinkService(storage)
	r := ls.Routes()

	short := ls.store.CreateLink(testedLongURL)

	req, err := http.NewRequest("GET", "/"+short, nil)
	if err != nil {
		t.Fatal(err)
	}

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusTemporaryRedirect
		p := w.Header().Get("Location")
		pageOK := err == nil && p == testedLongURL

		return statusOK && pageOK
	})
}

func TestLinkServer_GetLinkHandlerNegative(t *testing.T) {
	storage := linkstore.NewLinkStore()
	ls := NewLinkService(storage)
	r := ls.Routes()

	req, err := http.NewRequest("GET", "/wrongURL", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusNotFound
		p := w.Header().Get("Location")
		pageOK := err == nil && p != testedLongURL

		return statusOK && pageOK
	})
}
