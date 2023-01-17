package service

import (
	"context"
	"github.com/evgen1067/anti-bruteforce/internal/repository"
)

type BlacklistService struct {
	ctx context.Context
	db  repository.ListRepo
}

func NewBlacklistService(db repository.ListRepo, ctx context.Context) *BlacklistService {
	return &BlacklistService{
		db:  db,
		ctx: ctx,
	}
}

func (b *BlacklistService) AddToBlacklist(address string) error {
	return b.db.AddToBlacklist(b.ctx, address)
}

func (b *BlacklistService) ExistsInBlacklist(address string) (bool, error) {
	return b.db.ExistsInBlacklist(b.ctx, address)
}

func (b *BlacklistService) DeleteFromBlacklist(address string) error {
	return b.db.DeleteFromBlacklist(b.ctx, address)
}
