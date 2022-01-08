package app

import (
	"encoding/json"
	"github.com/magmel48/go-web/internal/auth"
	authmocks "github.com/magmel48/go-web/internal/auth/mocks"
	dbmocks "github.com/magmel48/go-web/internal/db/mocks"
	"github.com/DATA-DOG/go-sqlmock"
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
