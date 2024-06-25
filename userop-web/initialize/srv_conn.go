package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"shop-api/userop-web/global"
	"shop-api/userop-web/proto"
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

	/** 连接【用户相关操作服务】 **/
	userOpConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserOpSrvInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatalf("[InitSrvConn] 连接【用户相关操作服务】失败: %v", err)
	}
	global.AddressClient = proto.NewAddressClient(userOpConn)
	global.MessageClient = proto.NewMessageClient(userOpConn)
	global.UserFavClient = proto.NewUserFavClient(userOpConn)
}
