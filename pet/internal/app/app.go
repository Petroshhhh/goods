package app

import (
	"log/slog"
	grpcapp "pet/internal/app/grpc"
	"pet/internal/services/goodsHand"
	"pet/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
) *App {
	//инициализировать хранилише
	storage, err := postgres.New()
	if err != nil {
		panic(err)
	}
	// init auth service(auth)

	goodsHandService := goodsHand.New(log, storage, storage, storage)

	grpcApp := grpcapp.New(log, goodsHandService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
