package app

import (
	"errors"
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApp_handlePost(t *testing.T) {
	type fields struct {
		shortener shortener.Shortener
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type want struct {
		contentType  string
		statusCode   int
		shortenedURL string
	}

	noBodyRequest, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
	malformedURLInBodyRequest, _ := http.NewRequest(
		http.MethodPost, "http://localhost:8080", strings.NewReader("test"))
	okURLInBodyRequest, _ := http.NewRequest(
		http.MethodPost, "http://localhost:8080", strings.NewReader("https://google.com"))

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "no body",
			fields: fields{shortener: shortener.NewShortener("http://localhost:8080")},
			args: args{
				w: httptest.NewRecorder(),
				r: noBodyRequest,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:   "malformed url",
			fields: fields{shortener: shortener.NewShortener("http://localhost:8080")},
			args: args{
				w: httptest.NewRecorder(),
				r: malformedURLInBodyRequest,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:   "happy path",
			fields: fields{shortener: shortener.NewShortener("http://localhost:8080")},
			args: args{
				w: httptest.NewRecorder(),
				r: okURLInBodyRequest,
			},
			want: want{
				contentType:  "text/plain; charset=utf-8",
				statusCode:   http.StatusCreated,
				shortenedURL: "http://localhost:8080/1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener: tt.fields.shortener,
			}

			h := http.HandlerFunc(app.HandleHTTPRequests)
			h.ServeHTTP(tt.args.w, tt.args.r)

			switch w := tt.args.w.(type) {
			case *httptest.ResponseRecorder:
				result := w.Result()

				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

				if result.StatusCode == http.StatusCreated {
					rawBody, err := ioutil.ReadAll(result.Body)

					assert.NoError(t, err, "read body error")
					assert.Equal(t, tt.want.shortenedURL, (string)(rawBody))
				}
			default:
				assert.Error(t, errors.New("wrong test setup"))
			}
		})
	}
}

func TestApp_handleGet(t *testing.T) {
	type fields struct {
		shortener shortener.Shortener
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type want struct {
		contentType  string
		statusCode   int
		shortenedURL string
	}

	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "no url found for fresh db in shortener",
			fields: fields{shortener: shortener.NewShortener("http://localhost:8080")},
			args: args{
				w: httptest.NewRecorder(),
				r: request,
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{
				shortener: tt.fields.shortener,
			}

			h := http.HandlerFunc(app.HandleHTTPRequests)
			h.ServeHTTP(tt.args.w, tt.args.r)

			switch w := tt.args.w.(type) {
			case *httptest.ResponseRecorder:
				result := w.Result()

				assert.Equal(t, tt.want.statusCode, result.StatusCode)
				assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			default:
				assert.Error(t, errors.New("wrong test setup"))
			}
		})
	}
}
