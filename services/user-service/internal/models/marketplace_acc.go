package models

import (
	userpb "github.com/kimvlry/sales-sync/shared/proto/user"
)

type MarketplaceType string

const (
	Avito MarketplaceType = "avito"
)

func MarketplaceTypeFromProto(protoType userpb.MarketplaceType) MarketplaceType {
	switch protoType {
	case userpb.MarketplaceType_AVITO:
		return Avito
	default:
		return "unknown"
	}
}

type MarketplaceAccount struct {
	ID              string
	UserID          string
	MarketplaceType MarketplaceType
	AccountID       string
	Credentials     map[string]string
}

func (t MarketplaceType) ToProto() userpb.MarketplaceType {
	switch t {
	case Avito:
		return userpb.MarketplaceType_AVITO
	default:
		return userpb.MarketplaceType_UNKNOWN
	}
}

func (a *MarketplaceAccount) ToProto() *userpb.MarketplaceAccount {
	return &userpb.MarketplaceAccount{
		Id:              a.ID,
		UserId:          a.UserID,
		MarketplaceType: a.MarketplaceType.ToProto(),
		AccountId:       a.AccountID,
		Credentials:     a.Credentials,
	}
}
