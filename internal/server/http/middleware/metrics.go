package middleware

import (
	"go.elastic.co/apm"
	"net/http"
)

func(md *Middleware) ApmMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := apm.DefaultTracer.StartTransaction(r.URL.String(), "request")
		ctx := apm.ContextWithTransaction(r.Context(), t)
		next.ServeHTTP(w, r.WithContext(ctx))
		t.Result = "Success"
		t.End()
	})
}