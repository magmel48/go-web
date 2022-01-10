package app

import (
	"github.com/magmel48/go-web/internal/auth"
	"github.com/valyala/fasthttp"
	"log"
)

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
func cookiesHandler(authenticator auth.Auth) func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if authenticator != nil {
				_, err := getUserID(ctx, authenticator)

				// sets cookie if it's not valid (empty or wrong encoded)
				if err != nil {
					log.Println("user session invalidation error", err)

					userID := auth.NewUserID()
					userToken, _ := authenticator.Encode(userID)

					cookie := fasthttp.Cookie{}
					cookie.SetKey("session")
					cookie.SetValue(string(userToken))
					cookie.SetPath("/")
					cookie.SetHTTPOnly(true)

					ctx.Response.Header.SetCookie(&cookie)

					// in case of first user request we also need to set request cookie here
					// to be able to get it further
					ctx.Request.Header.SetCookie(string(cookie.Key()), string(cookie.Value()))
				}
			}

			h(ctx)
		}
	}
}

func getUserID(ctx *fasthttp.RequestCtx, authenticator auth.Auth) (auth.UserID, error) {
	sessionCookie := ctx.Request.Header.Cookie("session")
	return authenticator.Decode(sessionCookie)
}
