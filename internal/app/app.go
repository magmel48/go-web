package app

import (
	"fmt"
	"github.com/magmel48/go-web/internal/shortener"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type App struct {
	shortener shortener.Shortener
}

func NewApp(host string, port string) App {
	return App{
		shortener: shortener.NewShortener(fmt.Sprintf("http://%s:%s", host, port)),
	}
}

// HandleHTTPRequests handles http requests.
func (app App) HandleHTTPRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.handlePost(w, r)
	case http.MethodGet:
		app.handleGet(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (app App) handlePost(w http.ResponseWriter, r *http.Request) {
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read request body", http.StatusBadRequest)
		return
	}

	body := (string)(rawBody)
	_, err = url.Parse(body)
	if err != nil {
		http.Error(w, "cannot parse url", http.StatusBadRequest)
		return
	}

	shortURL := app.shortener.MakeShorter(body)

	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write(([]byte)(shortURL))
	if err != nil {
		http.Error(w, "cannot write response", http.StatusBadRequest)
	}
}

func (app App) handleGet(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 2 {
		http.Error(w, "cannot parse url", http.StatusBadRequest)
		return
	}

	id := path[len(path)-1]
	initialURL, err := app.shortener.RestoreLong(id)
	if err != nil {
		http.Error(w, "initial version of the link is not found", http.StatusBadRequest)
	}

	w.Header().Set("Location", initialURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

	_, err = w.Write([]byte{})
	if err != nil {
		http.Error(w, "cannot write response", http.StatusBadRequest)
	}
}
