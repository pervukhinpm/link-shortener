// Package service предоставляет сервисные компоненты для работы с данными.
// Включает в себя:
//   - Проверку состояния базы данных
//   - Работу с URL
package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PingService предоставляет функциональность для проверки состояния базы данных.
// Используется для мониторинга доступности хранилища данных.
type PingService struct {
	db *pgxpool.Pool
}

// NewPingService создает новый экземпляр PingService с указанным репозиторием.
// Инициализирует сервис для проверки состояния базы данных.
func NewPingService(db *pgxpool.Pool) *PingService {
	return &PingService{db: db}
}

// PingDB проверяет доступность базы данных.
// В случае успеха возвращает nil, в случае ошибки - описание проблемы.
func (p *PingService) PingDB(ctx context.Context) error {
	return p.db.Ping(ctx)
}
