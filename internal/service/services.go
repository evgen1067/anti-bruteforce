package service

import (
	"context"
	"github.com/evgen1067/anti-bruteforce/internal/bucket"
	"github.com/evgen1067/anti-bruteforce/internal/common"
	"github.com/evgen1067/anti-bruteforce/internal/repository/psql"
)

type Auth interface {
	Authorize(req common.APIAuthRequest) bool
}

type Blacklist interface {
	AddToBlacklist(address string) error
	ExistsInBlacklist(address string) (bool, error)
	DeleteFromBlacklist(address string) error
}

type Whitelist interface {
	AddToWhitelist(address string) error
	ExistsInWhitelist(address string) (bool, error)
	DeleteFromWhitelist(address string) error
}

type LeakyBucket interface {
	Add(req common.APIAuthRequest) bool
	ResetBucket()
}

type Services struct {
	Auth
	Blacklist
	Whitelist
}

func NewServices(ctx context.Context, db *psql.Repo, lb *bucket.LeakyBucket) *Services {
	blacklist := NewBlacklistService(db, ctx)
	whitelist := NewWhitelistService(db, ctx)
	auth := NewAuthService(ctx, blacklist, whitelist, lb)

	return &Services{
		Auth:      auth,
		Blacklist: blacklist,
		Whitelist: whitelist,
	}
}
