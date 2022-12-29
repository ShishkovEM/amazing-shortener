package linkserver

import (
	"bytes"
	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const testedLongURL = "http://ya.ru/"

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func TestLinkServer_CreateLinkHandler(t *testing.T) {
	r := gin.Default()
	ls := New()

	r.POST("/", ls.CreateLinkHandler)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(testedLongURL)))
	if err != nil {
		t.Fatal(err)
	}

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusCreated
		p, err := io.ReadAll(w.Body)
		pageOK := err == nil && strings.Contains(string(p), "http://localhost:8080/")

		return statusOK && pageOK
	})
}

func TestLinkServer_GetLinkHandlerPositive(t *testing.T) {
	r := gin.Default()
	ls := New()

	short := ls.store.CreateLink(testedLongURL)

	r.GET("/"+short, ls.GetLinkHandler)

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
	r := gin.Default()
	ls := New()

	short := ls.store.CreateLink(testedLongURL)

	r.GET("/"+short, ls.GetLinkHandler)

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

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *LinkServer
	}{
		{name: "new test #1", want: &LinkServer{store: linkstore.New()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
