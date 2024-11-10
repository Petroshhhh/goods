package goodsHand

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"pet/internal/domain/models"
	"pet/internal/storage"
)

// то, что он делает по настоящему
type GoodsHand struct {
	log         *slog.Logger
	gdsSaver    GoodsSaver
	gdsProvider GoodsProvider
	gdsDelete   GoodsDelete
	//gdsUpdate   GoodsUpdate
	//gdsList     GoodsList
}

type GoodsSaver interface {
	SaveGoods(
		ctx context.Context,
		brand string,
		placeSave int64,
		storeHouse int64,
		worker string,
	) (uid int64, err error)
}

type GoodsProvider interface {
	Goods(ctx context.Context, id int64) (models.Goods, error)
}

type GoodsDelete interface {
	DeleteGoods(ctx context.Context, id int64) (bool, error)
}

type GoodsUpdate interface {
	UpdateGoods(ctx context.Context, id int64) (bool, error)
}

type GoodsList interface {
	ListGoods(ctx context.Context,
		Brand string,
		PlaceSave int64,
		StoreHouse int64,
		Worker string,
		createdAfter string,
		pageNumber int64,
		pageSize int64,
	) ([]models.Goods, int64, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrGoodsExists        = errors.New("goods already exists")
	ErrGoodsNotFound      = errors.New("goods not found")
)

func New(
	log *slog.Logger,
	goodsSaver GoodsSaver,
	goodsProvider GoodsProvider,
	goodsDelete GoodsDelete,
	//	goodsUpdate GoodsUpdate,
	//	goodsList GoodsList,

) *GoodsHand {
	return &GoodsHand{
		log:         log,
		gdsSaver:    goodsSaver,
		gdsProvider: goodsProvider,
		gdsDelete:   goodsDelete,
		//gdsUpdate:   goodsUpdate,
		//gdsList:     goodsList,
	}
}

func (g *GoodsHand) Create(ctx context.Context,
	brand string,
	placeSave int64,
	storeHouse int64,
	worker string,
) (int64, error) {
	const op = "GoodsHand.Create"

	log := g.log.With(
		slog.String("op", op),
		slog.String("Brand", brand),
	)

	log.Info("attempting to create user")

	goods, err := g.gdsSaver.SaveGoods(ctx, brand, placeSave, storeHouse, worker)

	if err != nil {
		if errors.Is(err, storage.ErrGoodsExists) {
			g.log.Warn("goods already exists", err)

			return int64(0), fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		g.log.Error("failed to save goods", err)
		return int64(0), fmt.Errorf("%s: %w", op, err)
	}

	return goods, nil
}

//func (g *GoodsHand) Update(ctx context.Context,
//	goodsId int64,
//) (bool, error) {
//	const op = "GoodsHand.Update"
//
//	log := g.log.With(
//		slog.String("op", op),
//		slog.Int64("goodsId", goodsId),
//	)
//
//	log.Info("attempting to update user")
//
//	success, err := g.gdsUpdate.UpdateGoods(ctx, goodsId)
//
//	if err != nil {
//		if errors.Is(err, storage.Errmove) {
//			g.log.Warn("move is not success", err)
//
//			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
//		}
//
//		g.log.Error("failed to update goods", err)
//		return false, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return success, err
//}

func (g *GoodsHand) Delete(ctx context.Context,
	goodsId int64,
) (bool, error) {
	const op = "GoodsHand.Delete"

	log := g.log.With(
		slog.String("op", op),
		slog.Int64("goodsId", goodsId),
	)

	log.Info("attempting to Delete user")

	success, err := g.gdsDelete.DeleteGoods(ctx, goodsId)

	if err != nil {
		if errors.Is(err, storage.Errmove) {
			g.log.Warn("move is not success", err)

			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		g.log.Error("failed to Delete goods", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return success, err

}

func (g *GoodsHand) GetGoods(ctx context.Context,
	goodsId int64,
) (string, int64, int64, string, error) {
	const op = "GoodsHand.GetGoods"

	log := g.log.With(
		slog.String("op", op),
		slog.Int64("goodsId", goodsId),
	)

	log.Info("attempting to GetGoods user")

	goods, err := g.gdsProvider.Goods(ctx, goodsId)

	if err != nil {
		if errors.Is(err, storage.ErrGoodsNotFound) {
			g.log.Warn("goods not found", err)

			return goods.Brand, goods.PlaceSave, goods.StoreHouse, goods.Worker, fmt.Errorf("%s: %w", op, ErrGoodsExists)
		}

		g.log.Error("failed to GetGoods goods", err)
		return goods.Brand, goods.PlaceSave, goods.StoreHouse, goods.Worker, fmt.Errorf("%s: %w", op, err)
	}

	return goods.Brand, goods.PlaceSave, goods.StoreHouse, goods.Worker, err
}

//func (g *GoodsHand) ListGoods(ctx context.Context,
//	Brand string,
//	PlaceSave int64,
//	StoreHouse int64,
//	Worker string,
//	createdAfter string,
//	pageNumber int64,
//	pageSize int64,
//) ([]models.Goods, int64, error) {
//	const op = "GoodsHand.ListGoods"
//
//	log := g.log.With(
//		slog.String("op", op),
//		slog.String("Brand", Brand),
//	)
//
//	log.Info("attempting to ListGoods user")
//
//	goods, totalCount, err := g.gdsList.ListGoods(ctx, Brand, PlaceSave, StoreHouse, Worker, createdAfter, pageNumber, pageSize)
//
//	if err != nil {
//		if errors.Is(err, storage.ErrGoodsNotFound) {
//			g.log.Warn("goods not found", err)
//
//			return goods, totalCount, fmt.Errorf("%s: %w", op, ErrGoodsExists)
//		}
//
//		g.log.Error("failed to GetGoods goods", err)
//		return goods, totalCount, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return goods, totalCount, err
//}
