package errs

import "errors"

// ErrURLNotFound представляет ошибку, возникающую при попытке найти несуществующий URL.
// Используется для обработки случаев, когда запрашиваемый сокращенный URL не найден в базе данных.
var ErrURLNotFound = errors.New("shortened URL not found")
