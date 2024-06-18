package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"shop-api/order-web/global"
	"shop-api/order-web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo

	/** 连接【商品服务】 **/
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【商品服务】失败: %v", err)
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)

	/** 连接【库存服务】 **/
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【库存服务】失败: %v", err)
	}
	global.InventSrvClient = proto.NewInventoryClient(inventoryConn)

	/** 连接【订单服务】 **/
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【订单服务】失败: %v", err)
	}

	global.OrderSrvClient = proto.NewOrderClient(orderConn)
}
