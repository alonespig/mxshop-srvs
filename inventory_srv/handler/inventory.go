package handler

import (
	"context"
	"fmt"
	"mxshop/global"
	"mxshop/model"
	"mxshop/proto"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/golang/protobuf/ptypes/empty"
	goredislib "github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (s *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "库存不存在")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// 库存归还
// 1. 订单超时归还
// 2. 订单创建失败
// 3. 手动归还
func (s *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	tx := global.DB.Begin()
	for _, good := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: good.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存不存在")
		}
		//扣减，会出现数据不一致的问题，锁，分布式锁
		inv.Stocks += good.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &empty.Empty{}, nil
}

// 扣减库存
func (s *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("172.27.49.67:6379"),
	})
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	tx := global.DB.Begin()
	for _, good := range req.GoodsInfo {
		var inv model.Inventory
		mutex := rs.NewMutex(fmt.Sprintf("goods_lock_%d", good.GoodsId))
		if err := mutex.Lock(); err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "获取锁失败")
		}
		if result := global.DB.Where(&model.Inventory{Goods: good.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存不存在")
		}
		if inv.Stocks < good.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减，会出现数据不一致的问题，锁，分布式锁
		inv.Stocks -= good.Num

		tx.Save(&inv)

		if ok, err := mutex.Unlock(); !ok || err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "释放锁失败")
		}

		// tx.Save(&inv)
		// if result := tx.Model(&model.Inventory{}).Where("goods = ? and version = ?", good.GoodsId, inv.Version).Select("Stocks", "Version").Updates(model.Inventory{
		// 	Stocks:  inv.Stocks,
		// 	Version: inv.Version + 1,
		// }); result.RowsAffected == 0 {
		// 	tx.Rollback()
		// 	return nil, status.Errorf(codes.Internal, "库存扣减失败")
		// } else {
		// 	break
		// }
		// }
	}
	tx.Commit()
	return &empty.Empty{}, nil
}

// 这是库存
func (s *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*empty.Empty, error) {
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)

	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)

	return &empty.Empty{}, nil
}
