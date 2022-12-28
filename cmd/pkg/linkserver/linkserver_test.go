package linkserver

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	serverPort           = "8080"
	testedLongURL        = "http://ya.ru/"
	testedInvalidURL1    = "ftp://ssss.com"
	testedInvalidURL2    = "http://1://2"
	testNotFoundShortURL = "ThisURLCanNotBeFound"
)

//// Link Структура записи информации о гиперссылках
//type Link struct {
//	ID       int    // Идентификатор гиперссылки
//	Original string // Исходная (длинная) ссылка
//	Short    string // Короткая ссылка
//}
//
//type CreateLinkRequest struct {
//	LongURL string
//}
//
//type CreateLinkResponse struct {
//	ShortURL string
//}
//
//func serverAddress() string {
//	return "http://localhost:" + serverPort
//}
//
//func createLink(t *testing.T, test string) string {
//	linkValue := CreateLinkRequest{test}
//	linkStr := linkValue.LongURL
//	reqBody := bytes.NewBuffer([]byte(linkStr))
//
//	resp, err := http.Post(serverAddress()+"/", "text/plain", reqBody)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Fatalf("expected status code=Created, got %v", resp.StatusCode)
//	}
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	respValue := CreateLinkResponse{strings.Trim(string(body), "http://localhost:8080/")}
//	return respValue.ShortURL
//}
//
//func enforceContentTextPlain(t *testing.T, resp *http.Response) {
//	t.Helper()
//	contentType, ok := resp.Header["Content-Type"]
//	if !ok {
//		t.Fatalf("expected response header to have Content-Type, got %v", resp.Header)
//	}
//
//	if len(contentType) < 1 {
//		t.Fatalf("expected non-empty Content-Type, got %v", resp.Header)
//	}
//
//	ct := strings.ToLower(contentType[0])
//
//	if ct != "text/plain" && ct != "text/plain; charset=utf-8" {
//		t.Errorf("expected Content-Type=text/plain, got %v", ct)
//	}
//}
//
//func getLinkByShort(t *testing.T, short string) Link {
//	t.Helper()
//
//	client := &http.Client{
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return http.ErrUseLastResponse
//		},
//	}
//	req, err := http.NewRequest("GET", serverAddress()+fmt.Sprintf("/%s", short), nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	enforceContentTextPlain(t, resp)
//	if resp.StatusCode != http.StatusTemporaryRedirect {
//		t.Fatalf("expected status code=307(Temporary Redirect), got %v", resp.StatusCode)
//	}
//
//	originalURL, err := resp.Location()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	respValue := Link{0, originalURL.Scheme + "://" + originalURL.Host + originalURL.Path, short}
//
//	return respValue
//}
//
//func expectLinkNotFound(t *testing.T, short string) {
//	t.Helper()
//
//	client := &http.Client{
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return http.ErrUseLastResponse
//		},
//	}
//
//	req, err := http.NewRequest("GET", serverAddress()+fmt.Sprintf("/%s", short), nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusNotFound {
//		t.Fatalf("expected status code=StatusNotFound, got %v", resp.StatusCode)
//	}
//}

//func TestCreateAndGet(t *testing.T) {
//
//	short := createLink(t, testedLongURL)
//
//	link := getLinkByShort(t, short)
//
//	if link.Original != testedLongURL {
//		t.Errorf("expected found link %q, got %q", testedLongURL, link.Original)
//	}
//
//	expectLinkNotFound(t, testNotFoundShortURL)
//}

func TestLinkServer_createLinkHandlerPositive(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        201,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := New()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(testedLongURL)))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ls.createLinkHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) == strings.TrimPrefix(string(resBody), "http://localhost:8080/") {
				t.Errorf("Expected body %s, got %s", "http://localhost:8080/{id}", w.Body.String())
			}
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestLinkServer_createLinkHandlerNegative(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := New()
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(testedInvalidURL1)))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ls.createLinkHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			defer res.Body.Close()
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestLinkServer_getLinkHandlerNegative(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := New()
			request := httptest.NewRequest(http.MethodGet, "/aaaaa", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ls.getLinkHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}
			defer res.Body.Close()
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
