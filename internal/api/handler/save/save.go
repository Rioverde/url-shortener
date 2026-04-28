package save

import (
	"log/slog"
	"net/http"

	resp "github.com/Rioverde/url-shortener/internal/api/response"
	"github.com/Rioverde/url-shortener/internal/domain"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required, url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// New creates a new handler for saving a URL
func New(log *slog.Logger, service *domain.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var op = "hadlers.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("url", r.URL.String()),
		)

		var req Request
		// Decode the request body
		err := render.DecodeJSON(r.Body, &req)
		// If there is an error, log it and return a 400 Bad Request response
		if err != nil {
			log.Error("failed to decode request", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		// If there is no error, log the request body
		log.Info("Request Body decoded", slog.Any("request", req))
		// Validate the request body
		if err := validator.New().Struct(req); err != nil {
			// check if the error is a validation error
			validateErr, ok := err.(validator.ValidationErrors)
			if !ok {
				log.Error("failed to validate request", "error", err)
			}
			log.Error("request validation failed", "error", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		v, err := service.Shorten(req.URL)
		if err != nil {
			log.Error("failed to shorten URL", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to shorten URL"))
			return
		}

		resp := Response{
			Response: resp.OK(),
			Alias:    v,
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}
}
