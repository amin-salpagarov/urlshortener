package redirect

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


type URLGetter interface {
    GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handlers.url.redirect.New"

        log = log.With(
            slog.String("op", op),
            slog.String("request_id", middleware.GetReqID(r.Context())),
        )

        // Роутер chi позволяет делать вот такие финты -
        // получать GET-параметры по их именам.
        // Имена определяются при добавлении хэндлера в роутер, это будет ниже.
        alias := chi.URLParam(r, "alias")
        if alias == "" {
            log.Info("alias is empty")

            render.JSON(w, r, resp.Error("not found"))

            return
        }

        // Находим URL по алиасу в БД
        resURL, err := urlGetter.GetURL(alias)
        if errors.Is(err, storage.ErrURLNotFound) {
            // Не нашли URL, сообщаем об этом клиенту
            log.Info("url not found", "alias", alias)

            render.JSON(w, r, resp.Error("not found"))

            return
        }
        if err != nil {
            // Не удалось осуществить поиск
            log.Error("failed to get url", sl.Err(err))

            render.JSON(w, r, resp.Error("internal error"))

            return
        }

        log.Info("got url", slog.String("url", resURL))

        // Делаем редирект на найденный URL
        http.Redirect(w, r, resURL, http.StatusFound)
    }
}