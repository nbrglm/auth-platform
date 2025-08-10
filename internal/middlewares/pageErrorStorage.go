package middlewares

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

const CtxPageErrorKey = "pageError"
const CtxPageErrorReturnButtonTextKey = "pageErrorReturnButtonText"
const CtxPageErrorReturnURLKey = "pageErrorReturnURL"

func PageErrorStorageMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if pageErrCookie, err := ctx.Request.Cookie(CtxPageErrorKey); err != nil {
			handleErr(ctx)
			ctx.Next()
			return
		} else {
			if pageErr, err := url.QueryUnescape(pageErrCookie.Value); err != nil {
				handleErr(ctx)
				ctx.Next()
				return
			} else {
				ctx.Set(CtxPageErrorKey, pageErr)
			}
		}

		if returnButtonTextCookie, err := ctx.Request.Cookie(CtxPageErrorReturnButtonTextKey); err != nil {
			ctx.Set(CtxPageErrorReturnButtonTextKey, "")
		} else {
			if returnButtonText, err := url.QueryUnescape(returnButtonTextCookie.Value); err != nil {
				ctx.Set(CtxPageErrorReturnButtonTextKey, "")
			} else {
				ctx.Set(CtxPageErrorReturnButtonTextKey, returnButtonText)
			}
		}

		if returnUrlCookie, err := ctx.Request.Cookie(CtxPageErrorReturnURLKey); err != nil {
			ctx.Set(CtxPageErrorReturnURLKey, "")
		} else {
			if returnUrl, err := url.QueryUnescape(returnUrlCookie.Value); err != nil {
				ctx.Set(CtxPageErrorReturnURLKey, "")
			} else {
				ctx.Set(CtxPageErrorReturnURLKey, returnUrl)
			}
		}

		// Clear error cookies after reading them
		ctx.SetCookie(CtxPageErrorKey, "", -1, "/", "", false, true)
		ctx.SetCookie(CtxPageErrorReturnButtonTextKey, "", -1, "/", "", false, true)
		ctx.SetCookie(CtxPageErrorReturnURLKey, "", -1, "/", "", false, true)

		ctx.Next()
	}
}

func handleErr(ctx *gin.Context) {
	ctx.Set(CtxPageErrorKey, "")
	ctx.Set(CtxPageErrorReturnButtonTextKey, "")
	ctx.Set(CtxPageErrorReturnURLKey, "")
}

func SetPageError(c *gin.Context, errorMsg, returnText, returnURL string) {
	// URL encode to handle special characters
	c.SetCookie(CtxPageErrorKey, url.QueryEscape(errorMsg), 300, "/", "", false, true)
	if returnText != "" {
		c.SetCookie(CtxPageErrorReturnButtonTextKey, url.QueryEscape(returnText), 300, "/", "", false, true)
	}
	if returnURL != "" {
		c.SetCookie(CtxPageErrorReturnURLKey, url.QueryEscape(returnURL), 300, "/", "", false, true)
	}
}
