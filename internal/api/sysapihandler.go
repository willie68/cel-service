package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/willie68/cel-service/internal/serror"
)

// SysAPIKey defining a handler for checking system id and api key
type SysAPIConfig struct {
	Apikey string
	// Skip particular requests from the handler
	SkipFunc func(r *http.Request) bool
}

// SysAPIHandler creates a new directly usable handler
func SysAPIHandler(cfg SysAPIConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip tracer
			if cfg.SkipFunc != nil && cfg.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			if cfg.Apikey != strings.ToLower(r.Header.Get(APIKeyHeaderKey)) {
				msg := "apikey not correct"
				apierr := serror.BadRequest(nil, "missing-header", msg)
				render.Status(r, apierr.Code)
				render.JSON(w, r, apierr)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

var (
	ContextKeyOffset = contextKey("offset")
	ContextKeyLimit  = contextKey("limit")
)

//Paginate is a middleware logic for populating the context with offset and limit values
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		offsetStr := request.URL.Query().Get("offset")
		limitStr := request.URL.Query().Get("limit")
		if offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil {
				msg := "type of offset string is not correct."
				apierr := serror.BadRequest(err, "wrong-type", msg)
				render.Status(request, apierr.Code)
				render.JSON(response, request, apierr)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyOffset, offset)
		} else {
			ctx = context.WithValue(ctx, ContextKeyOffset, 0)
		}
		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				msg := "type of limit string is not correct."
				apierr := serror.BadRequest(err, "wrong-type", msg)
				render.Status(request, apierr.Code)
				render.JSON(response, request, apierr)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyLimit, limit)
		} else {
			ctx = context.WithValue(ctx, ContextKeyLimit, 0)
		}
		next.ServeHTTP(response, request.WithContext(ctx))
	})
}

type contextKey string

func (c contextKey) String() string {
	return "api" + string(c)
}
