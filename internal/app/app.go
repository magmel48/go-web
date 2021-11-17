package app

import (
	"fmt"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/valyala/fasthttp"
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
func (app App) HandleHTTPRequests(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case http.MethodPost:
		app.handlePost(ctx)
	case http.MethodGet:
		app.handleGet(ctx)
	default:
		ctx.NotFound()
	}
}

func (app App) handlePost(ctx *fasthttp.RequestCtx) {
	if ctx.PostBody() == nil {
		ctx.Error("empty request body", fasthttp.StatusBadRequest)
		return
	}

	body := string(ctx.PostBody())
	_, err := url.ParseRequestURI(body)
	if err != nil {
		ctx.Error( "cannot parse url", fasthttp.StatusBadRequest)
		return
	}

	shortURL := app.shortener.MakeShorter(body)

	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody([]byte(shortURL))
}

func (app App) handleGet(ctx *fasthttp.RequestCtx) {
	path := strings.Split(string(ctx.Path()), "/")
	if len(path) != 2 {
		ctx.Error("cannot parse url", fasthttp.StatusBadRequest)
		return
	}

	id := path[len(path)-1]
	initialURL, err := app.shortener.RestoreLong(id)
	if err != nil {
		ctx.Error("initial version of the link is not found", fasthttp.StatusBadRequest)
		return
	}

	ctx.Response.Header.Set("Location", initialURL)
	ctx.SetStatusCode(http.StatusTemporaryRedirect)
}
