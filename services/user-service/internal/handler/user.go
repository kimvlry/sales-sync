package handler

import (
	"context"
	"log/slog"

	userpb "github.com/kimvlry/sales-sync/shared/proto/user"
	"user-service/internal/logger"
	"user-service/internal/models"
	"user-service/internal/repository/postgres"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	repo postgres.UserRepository
}

func NewUserHandler(repo *postgres.UserRepository) *UserHandler {
	return &UserHandler{
		repo: *repo,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	ctx = logger.NewBuilder(ctx).WithTelegramID(req.TelegramId).Build()
	slog.InfoContext(ctx, "CreateUser called")

	user, err := h.repo.CreateUser(ctx, req.TelegramId, req.Name)
	if err != nil {
		slog.ErrorContext(ctx, "CreateUser failed", "error", err)
		return nil, err
	}
	ctx = logger.NewBuilder(ctx).WithUserID(user.ID).Build()
	slog.InfoContext(ctx, "CreateUser succeed")
	return &userpb.CreateUserResponse{
		User: user.ToProto(),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	ctx = logger.NewBuilder(ctx).WithTelegramID(req.Id).Build()
	slog.InfoContext(ctx, "GetUser called")

	user, err := h.repo.GetUser(ctx, req.GetId())
	if err != nil {
		slog.ErrorContext(ctx, "GetUser failed", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "GetUser succeed")
	return &userpb.GetUserResponse{
		User: user.ToProto(),
	}, nil
}

func (h *UserHandler) ConnectMarketplace(ctx context.Context, req *userpb.ConnectMarketplaceRequest) (*userpb.ConnectMarketplaceResponse, error) {
	ctx = logger.NewBuilder(ctx).
		WithUserID(req.UserId).
		WithExtra(map[string]any{
			"marketplace": req.Account.MarketplaceType.String(),
			"account_id":  req.Account.Id,
		}).Build()

	slog.InfoContext(ctx, "ConnectMarketplace called")

	account := models.MarketplaceAccount{
		MarketplaceType: models.MarketplaceTypeFromProto(req.Account.MarketplaceType),
		AccountID:       req.Account.AccountId,
		Credentials:     req.Account.Credentials,
	}

	err := h.repo.CreateMarketplaceAccount(ctx, account, req.UserId)
	if err != nil {
		slog.ErrorContext(ctx, "ConnectMarketplace failed", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "ConnectMarketplace succeed")
	return &userpb.ConnectMarketplaceResponse{
		Success: true,
	}, nil
}
