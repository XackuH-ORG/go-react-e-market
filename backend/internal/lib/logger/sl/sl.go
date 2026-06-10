package sl

import "log/slog"

// Возвращает атрибут для логирования ошибки.
// Пример использования:
//
//	err := someFunction()
//
//	if err != nil {
//	    logger.Error("An error occurred", sl.Err(err))
//	}
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
