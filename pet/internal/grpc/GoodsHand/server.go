package GoodsHand

import (
	"context"
	petv1 "github.com/Petroshhhh/protopet/gen/go/pet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GoodsHand interface {
	Create(
		ctx context.Context,
		Brand string,
		PlaceSave int64,
		StoreHouse int64,
		Worker string,
	) (goodsId int64, err error)
	//Update(
	//	ctx context.Context,
	//	goodsId int64,
	//) (success bool, err error)
	Delete(
		ctx context.Context,
		goodsId int64,
	) (success bool, err error)
	GetGoods(
		ctx context.Context,
		goodsId int64,
	) (brand string, placeSave int64, storeHouse int64, worker string, err error)
	//ListGoods(
	//	ctx context.Context,
	//	Brand string,
	//	PlaceSave int64,
	//	StoreHouse int64,
	//	Worker string,
	//	createdAfter string,
	//	pageNumber int64,
	//	pageSize int64,
	//) (list []*petv1.CreateRequest, totalCount int64, err error)
	//HistoryGoods(
	//	ctx context.Context,
	//	goodsId int64,
	//) (history []*petv1.ChangeHistory,
	//	err error)
	//// TODO:дописать. скорее всего надо делать через массив или карту.но пока хз

}

type serverAPI struct {
	petv1.UnimplementedGoodsHandServer //забрал структуру сервера!!!!!!!!!!!!!!!!
	goodsHand                          GoodsHand
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, goodsHand GoodsHand) {
	petv1.RegisterGoodsHandServer(gRPC, &serverAPI{goodsHand: goodsHand})
}

func validateCreate(req *petv1.CreateRequest) error {
	if req.GetBrand() == "" {
		return status.Error(codes.InvalidArgument, "missing brand")
	}

	if req.GetPlaceSave() == emptyValue {
		return status.Error(codes.InvalidArgument, "missing place save")
	}

	if req.GetStoreHouse() == emptyValue {
		return status.Error(codes.InvalidArgument, "missing store house")
	}
	if req.GetWorker() == "" {
		return status.Error(codes.InvalidArgument, "missing worker")
	}

	return nil
}

//func validateUpdate(req *petv1.UpdateRequest) error {
//
//	if req.GetGoodsId() == emptyValue {
//		return status.Error(codes.InvalidArgument, "missing goods id")
//	}
//
//	return nil
//}

func validateDelete(req *petv1.DeleteRequest) error {
	if req.GetGoodsId() == emptyValue {
		return status.Error(codes.InvalidArgument, "missing goods id")
	}

	return nil
}

func validateGetGoods(req *petv1.GetGoodsRequest) error {
	if req.GetGoodsId() == emptyValue {
		return status.Error(codes.InvalidArgument, "missing goods id")
	}

	return nil
}

//func validateListGoods(req *petv1.ListGoodsRequest) error {
//	if req.GetBrand() == "" {
//		return status.Error(codes.InvalidArgument, "missing brand")
//	}
//
//	if req.GetPlaceSave() == emptyValue {
//		return status.Error(codes.InvalidArgument, "missing place save")
//	}
//
//	if req.GetStoreHouse() == emptyValue {
//		return status.Error(codes.InvalidArgument, "missing store house")
//	}
//	if req.GetWorker() == "" {
//		return status.Error(codes.InvalidArgument, "missing worker")
//	}
//
//	if req.GetPageNumber() == emptyValue {
//		return status.Error(codes.InvalidArgument, "missing page number")
//	}
//
//	if req.GetPageSize() == emptyValue {
//		return status.Error(codes.InvalidArgument, "missing page size")
//	}
//
//	return nil
//}
//
//func validateHistoryGoods(req *petv1.HistoryGoodsRequest) error {
//	if req.GetGoodsId() == emptyValue {
//		return status.Error(codes.InvalidArgument, "missing goods id")
//	}
//
//	return nil
//}

func (s *serverAPI) Create(ctx context.Context,
	req *petv1.CreateRequest,
) (*petv1.CreateResponse, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}
	goods, err := s.goodsHand.Create(ctx, req.GetBrand(), req.GetPlaceSave(), req.GetStoreHouse(), req.GetWorker())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &petv1.CreateResponse{
		GoodsId: goods,
	}, nil
}

//func (s *serverAPI) Update(ctx context.Context,
//	req *petv1.UpdateRequest,
//) (*petv1.UpdateResponse, error) {
//	if err := validateUpdate(req); err != nil {
//		return nil, err
//	}
//	success, err := s.goodsHand.Update(ctx, req.GetGoodsId())
//	if err != nil {
//		return nil, status.Error(codes.Internal, "internal error")
//	}
//
//	return &petv1.UpdateResponse{
//		Success: success,
//	}, nil
//}

func (s *serverAPI) Delete(ctx context.Context,
	req *petv1.DeleteRequest,
) (*petv1.DeleteResponse, error) {
	if err := validateDelete(req); err != nil {
		return nil, err
	}
	success, err := s.goodsHand.Delete(ctx, req.GetGoodsId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &petv1.DeleteResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) GetGoods(ctx context.Context,
	req *petv1.GetGoodsRequest,
) (*petv1.GetGoodsResponse, error) {
	if err := validateGetGoods(req); err != nil {
		return nil, err
	}

	brand, placeSave, storeHouse, worker, err := s.goodsHand.GetGoods(ctx, req.GetGoodsId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &petv1.GetGoodsResponse{
		Brand:      brand,
		PlaceSave:  placeSave,
		StoreHouse: storeHouse,
		Worker:     worker,
	}, nil
}

//func (s *serverAPI) ListGoods(ctx context.Context,
//	req *petv1.ListGoodsRequest,
//) (*petv1.ListGoodsResponse, error) {
//	if err := validateListGoods(req); err != nil {
//		return nil, err
//	}
//
//	list, totalCount, err := s.goodsHand.ListGoods(ctx, req.GetBrand(), req.GetPlaceSave(), req.GetStoreHouse(), req.GetWorker(), req.GetCreatedAfter(), req.GetPageNumber(), req.GetPageSize())
//	if err != nil {
//		return nil, status.Error(codes.Internal, "internal error")
//	}
//
//	return &petv1.ListGoodsResponse{
//		Goods:      list,
//		TotalCount: totalCount,
//	}, nil
//}
//
//func (s *serverAPI) HistoryGoods(ctx context.Context,
//	req *petv1.HistoryGoodsRequest,
//) (*petv1.HistoryGoodsResponse, error) {
//	if err := validateHistoryGoods(req); err != nil {
//		return nil, err
//	}
//
//	list, err := s.goodsHand.HistoryGoods(ctx, req.GetGoodsId())
//	if err != nil {
//		return nil, status.Error(codes.Internal, "internal error")
//	}
//
//	return &petv1.HistoryGoodsResponse{
//		History: list,
//	}, nil
//}
