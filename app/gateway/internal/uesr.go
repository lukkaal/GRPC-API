package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukkaal/GRPC-API/app/gateway/rpc"
	userpb "github.com/lukkaal/GRPC-API/idl/user"
	"github.com/lukkaal/GRPC-API/pkg/errcode"
	"github.com/lukkaal/GRPC-API/pkg/jsonres"
	"github.com/lukkaal/GRPC-API/pkg/utils/jwt"
)

func UserRegister(ctx *gin.Context) {

	var userReq userpb.RegisterRequest
	if err := ctx.ShouldBind(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest,
			jsonres.RespErr(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.UserRegister(ctx, &userReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "UserRegister RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, jsonres.RespSuccess(ctx, r))
}

// userResp & token
func UserLogin(ctx *gin.Context) {
	var userReq userpb.LoginRequest
	if err := ctx.ShouldBind(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest,
			jsonres.RespErr(ctx, err, "绑定参数错误"))
		return
	}

	userResp, err := rpc.UserLogin(ctx, &userReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "加密错误"))
		return
	}

	// generate token
	token, err := jwt.GenerateToken(userResp.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "加密错误"))
		return
	}

	ctx.JSON(http.StatusOK, jsonres.RespSuccess(ctx,
		jsonres.TokenData{User: userResp, Token: token}))
}

func UserLogout(ctx *gin.Context) {
	// 解析 JWT，获取 userId
	// task 验证用户 token
	token := ctx.GetHeader("Authorization")

	if token == "" {
		code := 404
		ctx.JSON(200, gin.H{
			"status": code,
			"msg":    errcode.GetMsg(code),
		})
		ctx.Abort()
	}

	// validate the token
	_, err := jwt.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized,
			jsonres.RespErr(ctx, err, "无效的token"))
		return
	}
	// 可选：把 token 加到 Redis 黑名单，设置过期时间 = token 剩余有效期
	// redis.Set("blacklist:"+token, 1, tokenTTL)

	ctx.JSON(http.StatusOK, jsonres.RespSuccess(ctx, "退出成功"))
}
