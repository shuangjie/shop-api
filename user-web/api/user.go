package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"shop-api/user-web/middlewares"
	"shop-api/user-web/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/response"
	"shop-api/user-web/proto"
)

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandlerGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
		}
	}
}

func HandlerValidatorError(ctx *gin.Context, err error) {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// validator.ValidationErrors类型错误则进行翻译
	ctx.JSON(http.StatusOK, gin.H{
		"msg": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(ctx *gin.Context) {
	// 连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]", "msg", err.Error())
	}
	// 调用用户服务
	userSrvClient := proto.NewUserClient(userConn)

	// 前端传递的参数
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageSizeInt, _ := strconv.Atoi(pageSize)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:       uint32(pnInt),
		PageSize: uint32(pageSizeInt),
	})

	if err != nil {
		zap.S().Errorw("[GetUserList] 获取用户列表失败", "msg", err.Error())
		HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {

		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.Birthday), 0)).Format("2006-01-02"),
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)

}

func PassWordLogin(ctx *gin.Context) {
	//表单验证
	passWordLoginForm := forms.PassWordLoginForm{}

	if err := ctx.ShouldBind(&passWordLoginForm); err != nil {
		HandlerValidatorError(ctx, err)
		return
	}

	// 连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]", "msg", err.Error())
	}

	// 调用用户服务
	userSrvClient := proto.NewUserClient(userConn)

	//登录
	if rsp, loginErr := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	}); loginErr != nil {
		fmt.Println("loginErr")
		if e, ok := status.FromError(loginErr); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			}
			return
		}
	} else {
		//密码校验
		if passRsp, passErr := userSrvClient.CheckUserPassword(context.Background(), &proto.CheckPasswordInfo{
			PassWord:          passWordLoginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登录失败",
			})
		} else {
			if passRsp.Success {
				//生成 token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 7 天
						Issuer:    "Tomato",                       //签名的发行者
					},
				}
				token, tokenErr := j.CreateToken(claims)
				if tokenErr != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"token":      token,
					"nick_name":  rsp.NickName,
					"expired_at": (time.Now().Unix() + 60*60*24*7) * 1000,
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"msg": "登录失败",
				})
			}

		}

	}
}
