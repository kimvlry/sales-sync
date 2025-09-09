package models

import (
	userpb "github.com/kimvlry/sales-sync/shared/proto/user"
	"time"
)

type User struct {
	ID         string
	TelegramID string
	Name       string
	CreatedAt  time.Time
	Accounts   []MarketplaceAccount
}

func (u *User) ToProto() *userpb.User {
	protoAccounts := make([]*userpb.MarketplaceAccount, len(u.Accounts))
	for i, account := range u.Accounts {
		protoAccounts[i] = account.ToProto()
	}
	return &userpb.User{
		Id:         u.ID,
		TelegramId: u.TelegramID,
		Name:       u.Name,
		CreatedAt:  u.CreatedAt.Format("2006-01-02 15:04:05"),
		Accounts:   protoAccounts,
	}
}
