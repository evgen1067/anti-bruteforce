package service

import (
	"context"
	"github.com/evgen1067/anti-bruteforce/internal/common"
)

type AuthService struct {
	ctx       context.Context
	blacklist Blacklist
	whitelist Whitelist
}

func NewAuthService(ctx context.Context,
	blacklist Blacklist,
	whitelist Whitelist,
) *AuthService {
	return &AuthService{
		ctx:       ctx,
		blacklist: blacklist,
		whitelist: whitelist,
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
	// TODO "добавить дырявые ведра"
	return false
}
