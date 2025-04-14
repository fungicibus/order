package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	HeaderRequestID             = "X-Request-ID"
	ContextRequestID contextKey = "requestID"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(HeaderRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), ContextRequestID, rid)

		w.Header().Set(HeaderRequestID, rid)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if rid, ok := ctx.Value(ContextRequestID).(string); ok {
		return rid
	}
	return ""
}
