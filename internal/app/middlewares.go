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
				userID, err := getUserID(ctx, authenticator)

				// sets cookie if it's not valid (empty or wrong encoded)
				if err != nil {
					log.Println("user session invalidation error", err)

					userID = auth.NewUserID()
					userToken, _ := authenticator.Encode(userID)

					cookie := fasthttp.Cookie{}
					cookie.SetKey("session")
					cookie.SetValue(string(userToken))

					ctx.Response.Header.SetCookie(&cookie)
				}
			}

			h(ctx)
		}
	}
}

func getUserID(ctx *fasthttp.RequestCtx, authenticator auth.Auth) (auth.UserID, error) {
	// FIXME authenticator null?
	sessionCookie := ctx.Request.Header.Cookie("session")
	return authenticator.Decode(sessionCookie)
}
