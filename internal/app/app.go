package app

import (
	"encoding/json"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/config"
	"github.com/magmel48/go-web/internal/db"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4"
	"os"
)

// App makes urls shorter.
type App struct {
	shortener     shortener.Shortener
	authenticator auth.Auth
}

// ShortenPayload represents payload of request to /api/shorten.
type ShortenPayload struct {
	URL string `json:"url"`
}

// ShortenResult represents response from /api/shorten.
type ShortenResult struct {
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
		shortener:     shortener.NewShortener(baseURL, fileBackup),
		authenticator: authenticator,
	}
}

// HTTPHandler handles http requests.
func (app App) HTTPHandler() func(ctx *fasthttp.RequestCtx) {
	router := gorouter.NewFastHTTPRouter()
	router.POST("/", app.handlePost)
	router.POST("/api/shorten", app.handleJSONPost)
	router.GET("/{id:[0-9]+}", app.handleGet)
	router.GET("/user/urls", app.handleUserGet)
	router.GET("/ping", app.handlePing)

	return cookiesHandler(app.authenticator)(
		decompressHandler( // only for reading request
			fasthttp.CompressHandlerBrotliLevel( // only for writing response
				router.HandleFastHTTP, fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)))
}

func (app App) handlePost(ctx *fasthttp.RequestCtx) {
	if ctx.Request.Body() == nil {
		ctx.Error("empty request body", fasthttp.StatusBadRequest)
		return
	}

	userID, _ := getUserID(ctx, app.authenticator)
	body := string(ctx.Request.Body())
	shortURL, err := app.shortener.MakeShorter(body, userID)

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody([]byte(shortURL))
}

func (app App) handleJSONPost(ctx *fasthttp.RequestCtx) {
	var payload ShortenPayload

	body := ctx.Request.Body()
	err := json.Unmarshal(body, &payload)
	if err != nil {
		ctx.Error("wrong payload format", fasthttp.StatusBadRequest)
		return
	}

	userID, _ := getUserID(ctx, app.authenticator)
	shortURL, err := app.shortener.MakeShorter(payload.URL, userID)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	result := ShortenResult{
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
		ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
	default:
		ctx.Error("wrong id param", fasthttp.StatusBadRequest)
	}
}

func (app App) handleUserGet(ctx *fasthttp.RequestCtx) {
	userID, _ := getUserID(ctx, app.authenticator)
	result := app.shortener.GetUserLinks(userID)

	if len(result) == 0 {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	} else {
		response, err := json.Marshal(result)
		if err != nil {
			ctx.Error("json marshal error", fasthttp.StatusBadRequest)
			return
		}

		ctx.SetContentType("application/json; charset=utf-8")
		ctx.SetBody(response)
	}
}

func (app App) handlePing(ctx *fasthttp.RequestCtx) {
	if !db.CheckConnection() {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
