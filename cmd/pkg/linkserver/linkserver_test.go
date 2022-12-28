package linkserver

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	serverPort           = "8080"
	testedLongURL        = "http://ya.ru/"
	testNotFoundShortURL = "ThisURLCanNotBeFound"
)

// Link Структура записи информации о гиперссылках
type Link struct {
	ID       int    // Идентификатор гиперссылки
	Original string // Исходная (длинная) ссылка
	Short    string // Короткая ссылка
}

type CreateLinkRequest struct {
	LongURL string
}

type CreateLinkResponse struct {
	ShortURL string
}

func serverAddress() string {
	return "http://localhost:" + serverPort
}

func createLink(t *testing.T, test string) string {
	linkValue := CreateLinkRequest{test}
	linkStr := linkValue.LongURL
	reqBody := bytes.NewBuffer([]byte(linkStr))

	resp, err := http.Post(serverAddress()+"/", "text/plain", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code=Created, got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	respValue := CreateLinkResponse{strings.Trim(string(body), "http://localhost:8080/")}
	return respValue.ShortURL
}

func enforceContentTextPlain(t *testing.T, resp *http.Response) {
	t.Helper()
	contentType, ok := resp.Header["Content-Type"]
	if !ok {
		t.Fatalf("expected response header to have Content-Type, got %v", resp.Header)
	}

	if len(contentType) < 1 {
		t.Fatalf("expected non-empty Content-Type, got %v", resp.Header)
	}

	ct := strings.ToLower(contentType[0])

	if ct != "text/plain" && ct != "text/plain; charset=utf-8" {
		t.Errorf("expected Content-Type=text/plain, got %v", ct)
	}
}

func getLinkByShort(t *testing.T, short string) Link {
	t.Helper()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", serverAddress()+fmt.Sprintf("/%s", short), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	enforceContentTextPlain(t, resp)
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("expected status code=307(Temporary Redirect), got %v", resp.StatusCode)
	}

	originalURL, err := resp.Location()
	if err != nil {
		t.Fatal(err)
	}

	respValue := Link{0, originalURL.Scheme + "://" + originalURL.Host + originalURL.Path, short}

	return respValue
}

func expectLinkNotFound(t *testing.T, short string) {
	t.Helper()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", serverAddress()+fmt.Sprintf("/%s", short), nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status code=StatusNotFound, got %v", resp.StatusCode)
	}
}

func TestCreateAndGetIndirect(t *testing.T) {

	short := createLink(t, testedLongURL)

	link := getLinkByShort(t, short)

	if link.Original != testedLongURL {
		t.Errorf("expected found link %q, got %q", testedLongURL, link.Original)
	}

	expectLinkNotFound(t, testNotFoundShortURL)
}
