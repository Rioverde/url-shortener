package save

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/Rioverde/url-shortener/internal/api/response"
	"github.com/Rioverde/url-shortener/internal/domain"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// validate is reused across all requests. validator.Validate is safe for
// concurrent use and stateless — no need to instantiate it per request.
var validate = validator.New()

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// New creates a new handler for saving a URL.
func New(log *slog.Logger, service *domain.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.New"

		// Enrich the logger with per-request fields. Every line below
		// will carry op, request_id, and the requested URL path.
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("url", r.URL.String()),
		)

		// Decode the JSON body into our Request struct.
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		// Debug-level so we don't spam Info on every successful request.
		log.Debug("request body decoded", slog.Any("request", req))

		// Validate the request struct against the `validate` tags.
		if err := validate.Struct(req); err != nil {
			// validator.ValidationErrors is the expected shape; anything
			// else means the validator itself failed and that is a 500.
			var validateErr validator.ValidationErrors
			if !errors.As(err, &validateErr) {
				log.Error("validator failed", "error", err)
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("internal error"))
				return
			}
			log.Error("request validation failed", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		// Call the domain service. It owns the business rules and
		// returns sentinel errors that we translate into HTTP statuses below.
		alias, err := service.Shorten(req.URL)
		if err != nil {
			// Empty-URL is a client mistake — 400, not 500.
			if errors.Is(err, domain.ErrEmptyURL) {
				log.Info("rejected empty url")
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("url cannot be empty"))
				return
			}
			log.Error("failed to shorten url", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("url shortened", slog.String("alias", alias))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
