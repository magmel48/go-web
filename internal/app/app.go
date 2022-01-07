package app

import (
	"encoding/json"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4"
	"github.com/vardius/gorouter/v4/context"
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
func NewApp(baseURL string) App {
	authenticator, err := auth.NewCustomAuth()
	if err != nil {
		panic(err)
	}

	return App{
		shortener:     shortener.NewShortener(baseURL),
		authenticator: authenticator,
	}
}

// HTTPHandler handles http requests.
func (app App) HTTPHandler() func(ctx *fasthttp.RequestCtx) {
	router := gorouter.NewFastHTTPRouter()
	router.POST("/", app.handlePost)
	router.POST("/api/shorten", app.handleJSONPost)
	router.POST("/api/shorten/batch", app.handleBatchPost)
	router.GET("/user/urls", app.handleUserGet)
	router.GET("/ping", app.handlePing)
	router.GET("/{id}", app.handleGet)

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
	shortURL, isDuplicated, err := app.shortener.MakeShorter(ctx, body, userID)

	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetBody([]byte(shortURL))

	if isDuplicated {
		ctx.SetStatusCode(fasthttp.StatusConflict)
	} else {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	}
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
	shortURL, isDuplicated, err := app.shortener.MakeShorter(ctx, payload.URL, userID)
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
	ctx.SetBody(response)

	if isDuplicated {
		ctx.SetStatusCode(fasthttp.StatusConflict)
	} else {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	}
}

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

func (app App) handleGet(ctx *fasthttp.RequestCtx) {
	params := ctx.UserValue("params").(context.Params)
	id := params.Value("id")

	initialURL, err := app.shortener.RestoreLong(ctx, id)
	if err != nil {
		ctx.Error("initial version of the link is not found", fasthttp.StatusBadRequest)
		return
	}

	ctx.Response.Header.Set("Location", initialURL)
	ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
}

func (app App) handleUserGet(ctx *fasthttp.RequestCtx) {
	userID, _ := getUserID(ctx, app.authenticator)
	result, err := app.shortener.GetUserLinks(ctx, userID)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
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

func (app App) handlePing(ctx *fasthttp.RequestCtx) {
	if !db.CheckConnection(ctx) {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
