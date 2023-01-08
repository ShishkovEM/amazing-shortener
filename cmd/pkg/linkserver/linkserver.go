package linkserver

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ShishkovEM/amazing-shortener/internal/app/linkstore"

	"github.com/gin-gonic/gin"
)

const (
	serverAddress = "http://localhost"
	serverPort    = "8080"
)

type LinkServer struct {
	router *gin.Engine
	store  *linkstore.LinkStore
}

func NewLinkServer(store *linkstore.LinkStore) *LinkServer {
	ls := &LinkServer{
		router: gin.Default(),
		store:  store,
	}
	ls.setupRouter()
	return ls
}

func (ls *LinkServer) setupRouter() {
	ls.router.POST("/", ls.CreateLinkHandler)
	ls.router.GET("/:short", ls.GetLinkHandler)
}

func (ls *LinkServer) Run() {
	host := strings.TrimPrefix(serverAddress, "http://")
	ls.router.Run(host + ":" + serverPort)
}

func (ls *LinkServer) CreateLinkHandler(c *gin.Context) {

	LongURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	if !isValidURL(string(LongURL)) {
		c.String(http.StatusBadRequest, "Invalid URL creation request handled. Input URL: "+string(LongURL))
		return
	}

	short := ls.store.CreateLink(string(LongURL))
	c.Data(http.StatusCreated, "text/plain; charset=utf-8", []byte(serverAddress+":"+serverPort+"/"+short))
}

func (ls *LinkServer) GetLinkHandler(c *gin.Context) {

	short := c.Param("short")

	link, err := ls.store.GetLink(short)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Redirect(http.StatusTemporaryRedirect, link.Original)
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
