package main

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kimvlry/sales-sync/shared/pkg"
	"github.com/kimvlry/sales-sync/shared/pkg/db/conn"
	userpb "github.com/kimvlry/sales-sync/shared/proto/user"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"
	"user-service/internal/handler"
	"user-service/internal/logger"
	"user-service/internal/repository/postgres"
)

func main() {
	ctx := context.Background()

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	slog.Debug("loading config from", "path", cfgPath)

	cfg, err := pkg.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger.Init(cfg.App.Env)

	pool, err := conn.NewPostgresPool(ctx, cfg.DB.URL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	slog.Debug("connected to postgres", "url", cfg.DB.URL)
	defer pool.Close()

	sqlDB, err := sql.Open("pgx", cfg.DB.URL)
	if err != nil {
		log.Fatalf("failed to open db for migrations: %v", err)
	}
	defer sqlDB.Close()

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		log.Fatal("MIGRATIONS_PATH environment variable not set")
	}
	slog.Debug("loading migrations from", "path", migrationsPath)

	if err := goose.Up(sqlDB, migrationsPath); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	slog.Debug("migrations applied successfully")

	userRepo := postgres.NewUserRepository(pool)
	h := handler.NewUserHandler(userRepo)

	port := cfg.App.Port

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	slog.Info("UserService gRPC server listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
