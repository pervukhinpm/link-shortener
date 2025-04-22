package api

import (
	"net/http"

	"github.com/pervukhinpm/link-shortener.git/internal/service"
)

// DatabaseHealthHandler предоставляет HTTP-обработчик для проверки состояния базы данных.
// Используется для мониторинга доступности и работоспособности хранилища данных.
type DatabaseHealthHandler struct {
	ping *service.PingService
}

// NewDatabaseHealthHandler создает новый экземпляр DatabaseHealthHandler с указанным сервисом.
// Инициализирует обработчик для проверки состояния базы данных.
func NewDatabaseHealthHandler(
	ping *service.PingService,
) *DatabaseHealthHandler {
	return &DatabaseHealthHandler{
		ping: ping,
	}
}

// PingDatabase обрабатывает GET-запросы для проверки состояния базы данных.
// В случае успеха возвращает статус 200 OK, при ошибке - 500 Internal Server Error.
func (h *DatabaseHealthHandler) PingDatabase(w http.ResponseWriter, r *http.Request) {
	err := h.ping.PingDB(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
