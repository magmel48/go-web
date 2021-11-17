package app

import (
	"github.com/magmel48/go-web/internal/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"net"
	"testing"
)

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

func acquireRequest(method string, url string, body string) *fasthttp.Request {
	request := fasthttp.AcquireRequest()
	request.Header.SetMethod(method)
	request.SetRequestURI(url)

	if body != "" {
		request.SetBody([]byte(body))
	}

	return request
}

func TestApp_handlePost(t *testing.T) {
	type fields struct {
		shortener shortener.Shortener
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

	malformedURLInBodyRequest := acquireRequest(fasthttp.MethodPost, "http://localhost:8080", "test")
	okURLInBodyRequest := acquireRequest(fasthttp.MethodPost, "http://localhost:8080", "https://google.com")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "malformed url",
			fields: fields{shortener: shortener.NewShortener("http://localhost:8080")},
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
			name:   "happy path",
			fields: fields{shortener: shortener.NewShortener("http://localhost:8080")},
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
				shortener: tt.fields.shortener,
			}

			err := serve(app.HandleHTTPRequests, tt.args.r, tt.args.w)
			assert.NoError(t, err, "POST request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
			assert.Equal(t, tt.want.contentType, string(tt.args.w.Header.Peek("Content-Type")))

			if tt.args.w.StatusCode() == fasthttp.StatusCreated {
				assert.Equal(t, tt.want.shortenedURL, string(tt.args.w.Body()))
			}
		})
	}
}

func TestApp_handleGet(t *testing.T) {
	type fields struct {
		shortener shortener.Shortener
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

	request := acquireRequest(fasthttp.MethodGet, "http://localhost:8080/1", "")

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
				shortener: tt.fields.shortener,
			}

			err := serve(app.HandleHTTPRequests, tt.args.r, tt.args.w)
			assert.NoError(t, err, "GET request error")

			assert.Equal(t, tt.want.statusCode, tt.args.w.StatusCode())
			assert.Equal(t, tt.want.contentType, string(tt.args.w.Header.Peek("Content-Type")))
		})
	}
}
