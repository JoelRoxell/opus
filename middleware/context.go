package middleware

import (
	"builder/repositories"
	"context"
	"net/http"
)

const CustomContextKey = "CUSTOM_CONTEXT"

type CustomContext struct {
	Db *repositories.MongoDBDataStore
}

// ContextHandler is used to map reusable components to ctx.
func ContextHandler(f http.HandlerFunc, singletons CustomContext) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			ctx := req.Context()
			ctx = context.WithValue(ctx, CustomContextKey, singletons)

			f(res, req.WithContext(ctx))
		}()
	}
}
