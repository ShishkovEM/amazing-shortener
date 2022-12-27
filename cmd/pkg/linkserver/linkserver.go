package linkserver

import (
	"fmt"
	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"
	"io"
	"log"
	"net/http"
	"strings"
)

type LinkServer struct {
	store *linkstore.LinkStore
}

func New() *LinkServer {
	store := linkstore.New()
	return &LinkServer{store: store}
}

func (ls *LinkServer) LinkHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		if req.Method == http.MethodPost {
			ls.createLinkHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method POST, got %v", req.Method), http.StatusMethodNotAllowed)
		}
	} else {
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) != 1 {
			http.Error(w, "expect /<id> in link handler", http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodGet {
			ls.getLinkHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ls *LinkServer) createLinkHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling link create at %s\n", req.URL.Path)

	inputLink, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inputLinkParts := strings.Split(string(inputLink), "://")
	if len(inputLinkParts) != 2 {
		http.Error(w, "bad url", http.StatusBadRequest)
		return
	}

	if inputLinkParts[0] != "http" && inputLinkParts[0] != "https" {
		http.Error(w, "bad url", http.StatusBadRequest)
		return
	}

	newLinkID := ls.store.CreateLink(string(inputLink))

	responseBody, err := ls.store.GetLink(newLinkID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + responseBody.Short))
}

func (ls *LinkServer) getLinkHandler(w http.ResponseWriter, req *http.Request) {
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
