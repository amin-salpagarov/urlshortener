package remove


import (
	"errors"
	"net/http"
	resp "github.com/amin-salpagarov/urlshortener/internal/lib/api/response"
	sl "github.com/amin-salpagarov/urlshortener/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/amin-salpagarov/urlshortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
)




type URLRemover interface {
	DeleteUrl(alias string) (error)
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		err := urlRemover.DeleteUrl(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("url deleted", "alias", alias)

		render.JSON(w, r, resp.OK())

	}
}
