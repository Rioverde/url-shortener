package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/Rioverde/url-shortener/internal/api/response"
	"github.com/Rioverde/url-shortener/internal/domain"
	"github.com/Rioverde/url-shortener/internal/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

func New(log *slog.Logger, service *domain.URLService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is required"))
			return
		}

		url, err := service.GetUrl(alias)
		if err != nil {
			if errors.Is(err, domain.ErrEmptyKey) {
				log.Info("rejected empty key")
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, resp.Error("key cannot be empty"))
				return
			}
			if errors.Is(err, repo.ErrKeyNotFound) {
				log.Info("rejected key not found")
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("key not found"))
				return
			}
			log.Error("failed to get url", "error", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
