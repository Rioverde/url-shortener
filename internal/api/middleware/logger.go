package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// New returns an HTTP middleware that logs each incoming request and its outcome.
// Pair it with chi's RequestID middleware upstream so the request_id field is
// populated; otherwise it will be empty.
func New(log *slog.Logger) func(http.Handler) http.Handler {
	// Tag every log line emitted by this middleware with its component name.
	// Done once on construction, not per request, to avoid redundant work.
	log = log.With(slog.String("component", "middleware/logger"))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Build a child logger with request-specific fields.
			// Every log derived from `entry` will carry these.
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			// Wrap the writer so we can read status code and bytes written
			// after the handler returns. The plain http.ResponseWriter does
			// not expose these.
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			// Defer the completion log so it fires even if the handler panics
			// (assuming a recover middleware sits above this one).
			defer func() {
				// If the handler did not call WriteHeader explicitly, ww.Status()
				// returns 0. net/http treats that as 200, so normalise it here.
				status := ww.Status()
				if status == 0 {
					status = http.StatusOK
				}
				entry.Info("request completed",
					slog.Int("status", status),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
