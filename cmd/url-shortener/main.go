package main

import (
	"github.com/amin-salpagarov/urlshortener/internal/config"
	"log/slog"
	"os"
    "github.com/amin-salpagarov/urlshortener/internal/lib/logger/sl"
    "github.com/amin-salpagarov/urlshortener/internal/storage/sqlite"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/amin-salpagarov/urlshortener/internal/http-server/handlers/url/redirect"
    "github.com/amin-salpagarov/urlshortener/internal/http-server/handlers/url/save"
    "github.com/amin-salpagarov/urlshortener/internal/http-server/handlers/url/remove"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
    cfg := config.MustLoad()

    log := setupLogger(cfg.Env)
    log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

    storage, err := storage.New(cfg.StoragePath)
    if err != nil {
        log.Error("failed to initialize storage", sl.Err(err))
        return
    }

    router := chi.NewRouter()  
  
    router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
    router.Use(middleware.Logger) // Логирование всех запросов
    router.Use(middleware.Recoverer)  // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
    router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

    router.Post("/", save.New(log, storage))
    router.Get("/{alias}", redirect.New(log, storage))
    router.Delete("/{alias}", remove.New(log, storage))

    log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
    log.Debug("logger debug mode enabled")
}


func setupLogger(env string) *slog.Logger {
    var log *slog.Logger

    switch env {
    case envLocal:
        log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envDev:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    case envProd:
        log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
    }

    return log
}