package api

import (
	"context"
	"fmt"
	"net/http"
	"shop-api/user-web/forms"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"

	"shop-api/user-web/global"
)

// GenerateSmsCode 生成随机验证码
func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(uint64(time.Now().UnixNano()))

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}

func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandlerValidatorError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey(global.ServerConfig.AliyunSmsInfo.RegionID, global.ServerConfig.AliyunSmsInfo.AccessKeyID, global.ServerConfig.AliyunSmsInfo.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	smsCode := GenerateSmsCode(6)

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = global.ServerConfig.AliyunSmsInfo.RegionID
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile                     //手机号
	request.QueryParams["SignName"] = global.ServerConfig.AliyunSmsInfo.SignName //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_199575060"                        //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}"          //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	//  fmt.Print(response)
	if err != nil {
		fmt.Print(err.Error())
	}

	// 保存到 redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Minute*5)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
