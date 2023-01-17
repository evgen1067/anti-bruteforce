package service

import (
	"context"
	"github.com/evgen1067/anti-bruteforce/internal/common"
)

type AuthService struct {
	ctx       context.Context
	blacklist Blacklist
	whitelist Whitelist
	bucket    LeakyBucket
}

func NewAuthService(ctx context.Context,
	blacklist Blacklist,
	whitelist Whitelist,
	bucket LeakyBucket,
) *AuthService {
	return &AuthService{
		ctx:       ctx,
		blacklist: blacklist,
		whitelist: whitelist,
		bucket:    bucket,
	}
}

func (a *AuthService) Authorize(req common.APIAuthRequest) bool {
	res, err := a.whitelist.ExistsInWhitelist(req.IP)
	if err != nil {
		return false
	}
	if res {
		return true
	}

	res, err = a.blacklist.ExistsInBlacklist(req.IP)
	if err != nil {
		return false
	}
	if res {
		return false
	}
	return a.bucket.Add(req)
}
