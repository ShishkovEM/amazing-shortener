package linkserver

import (
	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

const (
	serverAddress = "http://localhost:"
	serverPort    = "8080"
)

type LinkServer struct {
	store *linkstore.LinkStore
}

func New() *LinkServer {
	store := linkstore.New()
	return &LinkServer{store: store}
}

func (ls *LinkServer) CreateLinkHandler(c *gin.Context) {

	LongURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	short := ls.store.CreateLink(string(LongURL))
	c.Data(http.StatusCreated, "text/plain; charset=utf-8", []byte(serverAddress+serverPort+"/"+short))
}

func (ls *LinkServer) GetLinkHandler(c *gin.Context) {

	link, err := ls.store.GetLink(strings.TrimPrefix(c.Request.URL.Path, "/"))
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Redirect(http.StatusTemporaryRedirect, link.Original)
}
