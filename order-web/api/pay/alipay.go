package pay

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"

	"shop-api/order-web/global"
	"shop-api/order-web/proto"
)

// client 支付宝客户端
func getClient(ctx *gin.Context) *alipay.Client {
	var err error
	var client *alipay.Client

	alipayInfo := global.ServerConfig.AlipayInfo

	if client, err = alipay.New(alipayInfo.AppID, alipayInfo.PrivateKey, alipayInfo.IsProduction); err != nil {
		zap.S().Errorf("[getClient] 创建 [支付宝客户端] 失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "支付宝客户端创建失败: " + err.Error(),
		})
	}
	if err = client.LoadAliPayPublicKey(alipayInfo.AlipayPublicKey); err != nil {
		zap.S().Errorf("[getClient] 加载支付宝公钥失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "加载支付宝公钥失败: " + err.Error(),
		})
	}

	return client
}

// Notify 支付宝 回调通知
func Notify(ctx *gin.Context) {
	client := getClient(ctx)
	// 解析表单数据
	err := ctx.Request.ParseForm()
	if err != nil {
		zap.S().Errorf("[Notify] 解析表单数据失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	notify, err := client.DecodeNotification(ctx.Request.Form)
	if err != nil {
		zap.S().Errorf("[Notify] 解析支付宝回调通知失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "解析支付宝回调通知失败: " + err.Error(),
		})
		return
	}

	// 业务逻辑处理
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: notify.OutTradeNo,
		Status:  string(notify.TradeStatus),
		TradeNo: notify.TradeNo,
	})
	if err != nil {
		zap.S().Errorf("[Notify] 更新订单状态失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "更新订单状态失败: " + err.Error(),
		})
		return
	}

	ctx.String(http.StatusOK, "success")

}

func GetPayUrl(ctx *gin.Context, orderSn string, orderAmount string) (string, error) {
	client := getClient(ctx)
	alipayInfo := global.ServerConfig.AlipayInfo

	var p = alipay.TradePagePay{}
	p.NotifyURL = alipayInfo.NotifyURL
	p.ReturnURL = alipayInfo.ReturnURL
	p.Subject = "支付测试:" + orderSn
	p.OutTradeNo = orderSn
	p.TotalAmount = orderAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorf("[GetPayUrl] 支付失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "支付失败: " + err.Error(),
		})
		return "", err
	}
	return url.String(), nil
}
