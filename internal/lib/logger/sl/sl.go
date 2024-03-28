package sl

import (
    "log/slog"

    //"github.com/amin-salpagarov/urlshortener/internal/lib/logger/handlers/slogdiscard"
)


func Err(err error) slog.Attr {
    return slog.Attr{
        Key:   "error",
        Value: slog.StringValue(err.Error()),
    }
}