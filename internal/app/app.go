package app

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/config"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
)

// App makes urls shorter.
type App struct {
	shortener shortener.Shortener
	auth      auth.Auth
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

	authenticator, err := auth.NewCustomAuth()
	if err != nil {
		panic(err)
	}

	return App{
		shortener: shortener.NewShortener(baseURL, fileBackup),
		auth: authenticator,
	}
}

// HTTPHandler handles http requests.
func (app App) HTTPHandler() func(ctx *fasthttp.RequestCtx) {
	router := fasthttprouter.New()
	router.POST("/", app.handlePost)
	router.POST("/api/shorten", app.handleJSONPost)
	router.GET("/:id", app.handleGet)

	return cookiesHandler(
		decompressHandler( // only for reading request
			fasthttp.CompressHandlerBrotliLevel( // only for writing response
				router.Handler, fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)))
}

// decompressHandler reads compressed request payload and decodes it.
func decompressHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		contentEncoding := ctx.Request.Header.Peek("Content-Encoding")
		switch string(contentEncoding) {
		case "gzip":
			body, err := ctx.Request.BodyGunzip()
			if err != nil {
				ctx.Error(err.Error(), fasthttp.StatusBadRequest)
				return
			}

			ctx.Request.SetBody(body)
		}

		h(ctx)
	}
}

// cookiesHandler sets and validates proper cookies.
func cookiesHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		sessionCookie := ctx.Request.Header.Cookie("session")
		if sessionCookie == nil {
			cookie := fasthttp.Cookie{}
			cookie.SetKey("session")
			//cookie.SetValue() // TODO get app somehow

			ctx.Response.Header.SetCookie(&cookie)
		} else {
			// FIXME check validity
		}

		h(ctx)
	}
}

func (app App) handlePost(ctx *fasthttp.RequestCtx) {
	if ctx.Request.Body() == nil {
		ctx.Error("empty request body", fasthttp.StatusBadRequest)
		return
	}

	body := string(ctx.Request.Body())
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

	body := ctx.Request.Body()
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
	fmt.Println(ctx.Path())

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

func (app App) handleUserGet(ctx *fasthttp.RequestCtx) {
	// TODO
}
