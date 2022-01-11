package app

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/magmel48/go-web/internal/auth"
	authmocks "github.com/magmel48/go-web/internal/auth/mocks"
	dbmocks "github.com/magmel48/go-web/internal/db/mocks"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"net"
	"testing"
)

var emptyHeaders = make(map[string]string)

// serve helps to run fasthttp mock server and send a request to created server.
func serve(handler fasthttp.RequestHandler, req *fasthttp.Request, res *fasthttp.Response) error {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(err)
		}
	}()

	client := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return ln.Dial()
		},
	}

	return client.Do(req, res)
}

func acquireRequest(method string, url string, body string, headers map[string]string) *fasthttp.Request {
	request := fasthttp.AcquireRequest()
	request.Header.SetMethod(method)
	request.SetRequestURI(url)

	if body != "" {
		request.SetBody([]byte(body))
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return request
}

func TestApp_handlePost(t *testing.T) {
	type fields struct {
		shortener     shortener.Shortener
		authenticator auth.Auth
	}
	type args struct {
		w *fasthttp.Response
		r *fasthttp.Request
	}
	type want struct {
		contentType  string
		statusCode   int
		shortenedURL string
	}

	mockAuth := &authmocks.Auth{}
	mockAuth.On("Decode", mock.Anything).Return(nil, nil)
	mockAuth.On("Encode", mock.Anything).Return(nil)

	mockDB := &dbmocks.DB{}
	mockDB.On("CreateSchema").Return(nil)

	db, sqlMock, _ := sqlmock.New()
	mockDB.On("Instance").Return(db)

	rows := sqlmock.NewRows([]string{"id", "short_id"}).AddRow(1, "1")
	sqlMock.ExpectQuery(`INSERT INTO "links"`).WillReturnRows(rows)

	malformedURLInBodyRequest := acquireRequest(
		fasthttp.MethodPost, "http://localhost:8080", "test", emptyHeaders)
	okURLInBodyRequest := acquireRequest(
		fasthttp.MethodPost, "http://localhost:8080", "https://google.com", emptyHeaders)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "malformed url",
			fields: fields{
				shortener:     shortener.NewShortener("http://localhost:8080", mockDB),
				authenticator: mockAuth,
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: malformedURLInBodyRequest,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  fasthttp.StatusBadRequest,
			},
		},
		{
			name: "happy path",
			fields: fields{
				shortener:     shortener.NewShortener("http://localhost:8080", mockDB),
				authenticator: mockAuth,
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: okURLInBodyRequest,
			},
			want: want{
				contentType:  "text/plain; charset=utf-8",
				statusCode:   fasthttp.StatusCreated,
				shortenedURL: "http://localhost:8080/1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener:     tt.fields.shortener,
				authenticator: tt.fields.authenticator,
			}

			err := serve(app.HTTPHandler(), tt.args.r, tt.args.w)
			assert.NoError(t, err, "POST request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
			assert.Equal(t, tt.want.contentType, string(tt.args.w.Header.Peek("Content-Type")))

			if tt.args.w.StatusCode() == fasthttp.StatusCreated {
				assert.Equal(t, tt.want.shortenedURL, string(tt.args.w.Body()))
			}
		})
	}
}

func TestApp_handleJSONPost(t *testing.T) {
	type fields struct {
		shortener     shortener.Shortener
		authenticator auth.Auth
	}
	type args struct {
		w *fasthttp.Response
		r *fasthttp.Request
	}
	type want struct {
		contentType string
		statusCode  int
		result      ShortenResult
	}

	mockAuth := &authmocks.Auth{}
	mockAuth.On("Decode", mock.Anything).Return(nil, nil)
	mockAuth.On("Encode", mock.Anything).Return(nil)

	mockDB := &dbmocks.DB{}
	mockDB.On("CreateSchema").Return(nil)

	db, sqlMock, _ := sqlmock.New()
	mockDB.On("Instance").Return(db)

	rows := sqlmock.NewRows([]string{"id", "short_id"}).AddRow(1, "1")
	sqlMock.ExpectQuery(`INSERT INTO "links"`).WillReturnRows(rows)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	malformedBody, _ := json.Marshal("[1,2,3]")
	okBody, _ := json.Marshal(ShortenPayload{URL: "https://google.com"})

	malformedURLInBodyRequest := acquireRequest(
		fasthttp.MethodPost, "http://localhost:8080/api/shorten", string(malformedBody), headers)
	okURLInBodyRequest := acquireRequest(
		fasthttp.MethodPost, "http://localhost:8080/api/shorten", string(okBody), headers)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "malformed url",
			fields: fields{
				shortener:     shortener.NewShortener("http://localhost:8080", mockDB),
				authenticator: mockAuth,
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: malformedURLInBodyRequest,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  fasthttp.StatusBadRequest,
			},
		},
		{
			name: "happy path",
			fields: fields{
				shortener:     shortener.NewShortener("http://localhost:8080", mockDB),
				authenticator: mockAuth,
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: okURLInBodyRequest,
			},
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  fasthttp.StatusCreated,
				result:      ShortenResult{Result: "http://localhost:8080/1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener:     tt.fields.shortener,
				authenticator: tt.fields.authenticator,
			}

			err := serve(app.HTTPHandler(), tt.args.r, tt.args.w)
			assert.NoError(t, err, "POST request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
			assert.Equal(t, tt.want.contentType, string(tt.args.w.Header.Peek("Content-Type")))

			if tt.args.w.StatusCode() == fasthttp.StatusCreated {
				var result ShortenResult

				err = json.Unmarshal(tt.args.w.Body(), &result)
				assert.NoError(t, err, "unmarshal response error")
				assert.Equal(t, tt.want.result, result)
			}
		})
	}
}

func TestApp_handleGet(t *testing.T) {
	type fields struct {
		shortener     shortener.Shortener
		authenticator auth.Auth
	}
	type args struct {
		w *fasthttp.Response
		r *fasthttp.Request
	}
	type want struct {
		contentType  string
		statusCode   int
		shortenedURL string
	}

	mockAuth := &authmocks.Auth{}
	mockAuth.On("Decode", mock.Anything).Return(nil, nil)
	mockAuth.On("Encode", mock.Anything).Return(nil)

	mockDB := &dbmocks.DB{}
	mockDB.On("CreateSchema").Return(nil)

	db, sqlMock, _ := sqlmock.New()
	mockDB.On("Instance").Return(db)

	rows := sqlmock.NewRows([]string{"id", "short_id", "original_url"})
	sqlMock.ExpectQuery(`SELECT "id", "short_id", "original_url" FROM "links"`).WillReturnRows(rows)

	request := acquireRequest(fasthttp.MethodGet, "http://localhost:8080/1", "", emptyHeaders)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "no url found for fresh db in shortener",
			fields: fields{
				shortener:     shortener.NewShortener("http://localhost:8080", mockDB),
				authenticator: mockAuth,
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: request,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  fasthttp.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener:     tt.fields.shortener,
				authenticator: tt.fields.authenticator,
			}

			err := serve(app.HTTPHandler(), tt.args.r, tt.args.w)
			assert.NoError(t, err, "GET request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
			assert.Equal(t, tt.want.contentType, string(tt.args.w.Header.Peek("Content-Type")))
		})
	}
}

func TestApp_handleUserGet(t *testing.T) {
	type fields struct {
		shortener     shortener.Shortener
		authenticator auth.Auth
	}
	type args struct {
		w *fasthttp.Response
		r *fasthttp.Request
	}
	type want struct {
		contentType string
		statusCode  int
		result      string
	}

	request := acquireRequest(fasthttp.MethodGet, "http://localhost:8080/user/urls", "", emptyHeaders)

	userID := "user_id_1"

	mockAuth := &authmocks.Auth{}
	mockAuth.On("Decode", mock.Anything).Return(&userID, nil)

	mockDB := &dbmocks.DB{}
	mockDB.On("CreateSchema").Return(nil)

	db, sqlMock, _ := sqlmock.New()
	mockDB.On("Instance").Return(db)

	rows := sqlmock.NewRows([]string{"short_id", "original_url"}).AddRow("1", "https://google.com")
	sqlMock.ExpectQuery(
		`SELECT l."short_id", l."original_url" FROM "user_links"`).WillReturnRows(rows)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "gets all user related urls",
			fields: fields{
				shortener:     shortener.NewShortener("http://localhost:8080", mockDB),
				authenticator: mockAuth,
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: request,
			},
			want: want{
				contentType: "application/json; charset=utf-8",
				statusCode:  fasthttp.StatusOK,
				result:      `[{"original_url":"https://google.com", "short_url":"http://localhost:8080/1"}]`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener:     tt.fields.shortener,
				authenticator: tt.fields.authenticator,
			}

			err := serve(app.HTTPHandler(), tt.args.r, tt.args.w)
			assert.NoError(t, err, "GET request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
			assert.Equal(t, tt.want.contentType, string(tt.args.w.Header.Peek("Content-Type")))
			assert.JSONEq(t, tt.want.result, string(tt.args.w.Body()))
		})
	}
}

func TestApp_handlePing(t *testing.T) {
	type fields struct {
		shortener     shortener.Shortener
		authenticator auth.Auth
	}
	type args struct {
		w *fasthttp.Response
		r *fasthttp.Request
	}
	type want struct {
		statusCode int
	}

	pingRequest := acquireRequest(
		fasthttp.MethodGet, "http://localhost:8080/ping", "", emptyHeaders)

	mockDB := &dbmocks.DB{}
	mockDB.On("CreateSchema").Return(nil)
	mockDB.On("Instance").Return(nil)
	mockDB.On("CheckConnection", mock.Anything).Return(false)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "returns 500 if connection is not ok",
			fields: fields{
				shortener: shortener.NewShortener("http://localhost:8080", mockDB),
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: pingRequest,
			},
			want: want{
				statusCode: fasthttp.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener:     tt.fields.shortener,
				authenticator: tt.fields.authenticator,
			}

			err := serve(app.HTTPHandler(), tt.args.r, tt.args.w)
			assert.NoError(t, err, "GET request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
		})
	}
}

func TestApp_handleDelete(t *testing.T) {
	type fields struct {
		shortener     shortener.Shortener
		authenticator auth.Auth
	}
	type args struct {
		w *fasthttp.Response
		r *fasthttp.Request
	}
	type want struct {
		statusCode int
	}

	request := acquireRequest(
		fasthttp.MethodDelete, "http://localhost:8080/api/user/urls", `["1"]`, emptyHeaders)

	mockAuth := &authmocks.Auth{}
	mockAuth.On("Decode", mock.Anything).Return(nil, nil)
	mockAuth.On("Encode", mock.Anything).Return(nil)

	mockDB := &dbmocks.DB{}
	mockDB.On("CreateSchema").Return(nil)
	mockDB.On("Instance").Return(nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "happy path",
			fields: fields{
				shortener: shortener.NewShortener("http://localhost:8080", mockDB),
				authenticator: mockAuth,
			},
			args: args{
				w: fasthttp.AcquireResponse(),
				r: request,
			},
			want: want{
				statusCode: fasthttp.StatusAccepted,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener:     tt.fields.shortener,
				authenticator: tt.fields.authenticator,
			}

			err := serve(app.HTTPHandler(), tt.args.r, tt.args.w)
			assert.NoError(t, err, "DELETE request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
		})
	}
}
