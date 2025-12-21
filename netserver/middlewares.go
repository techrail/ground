package netserver

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/techrail/ground/constants/customCtxKey"
)

type middleware struct{}

var Middleware *middleware

func init() {
	Middleware = new(middleware)
}

// RequestIDMiddleware is a Middleware to check and set requestID
func (m *middleware) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(customCtxKey.RequestId)
		if requestID == "" {
			// Generate a new UUID as the requestID
			requestID = uuid.New().String()
		}
		// Create a new context with the userID value
		ctx := context.WithValue(r.Context(), customCtxKey.RequestId, requestID)

		// Create a new request with the new context
		r = r.WithContext(ctx)

		// Call the next handler in the chain with the updated request
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.Printf("Panic recovered: %v", err)
				// Respond with 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// File ends here
