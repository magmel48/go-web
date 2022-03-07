package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db"
	"github.com/magmel48/go-web/internal/db/links"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4"
	routercontext "github.com/vardius/gorouter/v4/context"
)

// App makes urls shorter.
type App struct {
	shortener     shortener.Shortener
	authenticator auth.Auth
}

// ShortenPayload represents payload of a request to /api/shorten.
type ShortenPayload struct {
	URL string `json:"url"`
}

// ShortenResult represents response from /api/shorten.
type ShortenResult struct {
	Result string `json:"result"`
}

// BatchPayloadElement is one element from array from payload of a request to /api/shorten/batch.
type BatchPayloadElement struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResultElement is one element from array from response from /api/shorten/batch.
type BatchResultElement struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// NewApp creates new app that handles requests for making url shorter.
func NewApp(ctx context.Context, baseURL string) App {
	authenticator, err := auth.NewCustomAuth()
	if err != nil {
		panic(err)
	}

	database := db.SQLDB{}
	if err := database.CreateSchema(); err != nil {
		panic(err)
	}

	return App{
		shortener:     shortener.NewShortener(ctx, baseURL, &database),
		authenticator: authenticator,
	}
}

// HTTPHandler handles http requests.
func (app App) HTTPHandler() func(ctx *fasthttp.RequestCtx) {
	router := gorouter.NewFastHTTPRouter()
	router.POST("/", app.handlePost)
	router.POST("/api/shorten", app.handleJSONPost)
	router.POST("/api/shorten/batch", app.handleBatchPost)
	router.GET("/api/user/urls", app.handleUserGet)
	router.GET("/ping", app.handlePing)
	router.GET("/{id}", app.handleGet)
	router.DELETE("/api/user/urls", app.handleDelete)

	return cookiesHandler(app.authenticator)(
		decompressHandler( // only for reading request
			fasthttp.CompressHandlerBrotliLevel( // only for writing response
				router.HandleFastHTTP, fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)))
}

// handlePost handles POST on "/" route - creates new short link.
func (app App) handlePost(ctx *fasthttp.RequestCtx) {
	if ctx.Request.Body() == nil {
		ctx.Error("empty request body", fasthttp.StatusBadRequest)
		return
	}

	userID, err := getUserID(ctx, app.authenticator)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	body := string(ctx.Request.Body())
	shortURL, err := app.shortener.MakeShorter(ctx, body, userID)

	if err != nil {
		if !errors.Is(err, links.ErrConflict) {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusConflict)
	} else {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	}

	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetBody([]byte(shortURL))
}

// handleJSONPost does the same as handlePost, but accepts JSON payload and available on "/api/shorten".
func (app App) handleJSONPost(ctx *fasthttp.RequestCtx) {
	var payload ShortenPayload

	body := ctx.Request.Body()
	err := json.Unmarshal(body, &payload)
	if err != nil {
		ctx.Error("wrong payload format", fasthttp.StatusBadRequest)
		return
	}

	userID, err := getUserID(ctx, app.authenticator)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	shortURL, err := app.shortener.MakeShorter(ctx, payload.URL, userID)
	if err != nil {
		if !errors.Is(err, links.ErrConflict) {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusConflict)
	} else {
		ctx.SetStatusCode(fasthttp.StatusCreated)
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
	ctx.SetBody(response)
}

// handleBatchPost accepts multiple JSON records for making shorter links, available on "/api/shorten/batch".
func (app App) handleBatchPost(ctx *fasthttp.RequestCtx) {
	var payload []BatchPayloadElement

	body := ctx.Request.Body()
	err := json.Unmarshal(body, &payload)
	if err != nil {
		ctx.Error("wrong payload format", fasthttp.StatusBadRequest)
		return
	}

	originalURLs := make([]string, len(payload))
	for i, el := range payload {
		originalURLs[i] = el.OriginalURL
	}

	shortURLs, err := app.shortener.MakeShorterBatch(ctx, originalURLs)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}

	result := make([]BatchResultElement, len(payload))
	for i, el := range shortURLs {
		result[i] = BatchResultElement{CorrelationID: payload[i].CorrelationID, ShortURL: el}
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

// handleGet handles GET on "/api/user/urls" and returns original link from specified identifier.
func (app App) handleGet(ctx *fasthttp.RequestCtx) {
	params := ctx.UserValue("params").(routercontext.Params)
	id := params.Value("id")

	initialURL, err := app.shortener.RestoreLong(ctx, id)
	if err != nil {
		if errors.Is(err, shortener.ErrDeleted) {
			ctx.SetStatusCode(fasthttp.StatusGone)
			return
		}

		ctx.Error("initial version of the link is not found", fasthttp.StatusBadRequest)
		return
	}

	ctx.Response.Header.Set("Location", initialURL)
	ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
}

// handleUserGet handles GET on "/api/user/urls" and returns all links from the user.
func (app App) handleUserGet(ctx *fasthttp.RequestCtx) {
	userID, err := getUserID(ctx, app.authenticator)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	result, err := app.shortener.GetUserLinks(ctx, userID)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

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

// handlePing handles GET on "/ping" and checks database availability.
func (app App) handlePing(ctx *fasthttp.RequestCtx) {
	if !app.shortener.IsStorageAvailable(ctx) {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

// handleDelete handles DELETE on "/api/user/urls" and asynchronously deletes specified links.
func (app App) handleDelete(ctx *fasthttp.RequestCtx) {
	userID, err := getUserID(ctx, app.authenticator)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}

	var payload []string

	body := ctx.Request.Body()
	err = json.Unmarshal(body, &payload)
	if err != nil {
		ctx.Error("wrong payload format", fasthttp.StatusBadRequest)
		return
	}

	go app.shortener.DeleteURLs(userID, payload)
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}
