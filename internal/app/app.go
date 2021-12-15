package app

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/magmel48/go-web/internal/config"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
)

// App makes urls shorter.
type App struct {
	shortener shortener.Shortener
}

// Payload represents payload of request to /api/shorten.
type Payload struct {
	URL string `json:"url"`
}

// Result represents response from /api/shorten.
type Result struct {
	Result string `json:"result"`
}

// NewApp creates new app that handles requests for making url shorter.
func NewApp(baseURL string) App {
	fileBackup, err := shortener.NewFileBackup(config.FilePath, os.OpenFile)
	if err != nil {
		panic(err)
	}

	return App{
		shortener: shortener.NewShortener(baseURL, fileBackup),
	}
}

// HTTPHandler handles http requests.
func (app App) HTTPHandler() func(ctx *fasthttp.RequestCtx) {
	router := fasthttprouter.New()
	router.POST("/", app.handlePost)
	router.POST("/api/shorten", app.handleJSONPost)
	router.GET("/:id", app.handleGet)

	return fasthttp.CompressHandlerBrotliLevel(
		router.Handler, fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)
}

func (app App) handlePost(ctx *fasthttp.RequestCtx) {
	if ctx.PostBody() == nil {
		ctx.Error("empty request body", fasthttp.StatusBadRequest)
		return
	}

	body := string(ctx.PostBody())
	shortURL, err := app.shortener.MakeShorter(body)

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody([]byte(shortURL))
}

func (app App) handleJSONPost(ctx *fasthttp.RequestCtx) {
	var payload Payload

	body := ctx.PostBody()
	err := json.Unmarshal(body, &payload)
	if err != nil {
		ctx.Error("wrong payload format", fasthttp.StatusBadRequest)
		return
	}

	shortURL, err := app.shortener.MakeShorter(payload.URL)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	result := Result{
		Result: shortURL,
	}

	response, err := json.Marshal(result)
	if err != nil {
		ctx.Error("json marshal error", fasthttp.StatusBadRequest)
		return
	}

	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(response)
}

func (app App) handleGet(ctx *fasthttp.RequestCtx) {
	rawID := ctx.UserValue("id")

	switch id := rawID.(type) {
	case string:
		initialURL, err := app.shortener.RestoreLong(id)
		if err != nil {
			ctx.Error("initial version of the link is not found", fasthttp.StatusBadRequest)
			return
		}

		ctx.Response.Header.Set("Location", initialURL)
		ctx.SetStatusCode(http.StatusTemporaryRedirect)
	default:
		ctx.Error("wrong id param", http.StatusBadRequest)
	}
}
