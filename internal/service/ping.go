package service

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PingService struct {
	db *pgxpool.Pool
}

func NewPingService(db *pgxpool.Pool) *PingService {
	return &PingService{db: db}
}

func (p *PingService) PingDB(ctx context.Context) error {
	return p.db.Ping(ctx)
}
