package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukkaal/GRPC-API/app/gateway/rpc"
	taskpb "github.com/lukkaal/GRPC-API/idl/task"
	"github.com/lukkaal/GRPC-API/pkg/jsonres"
)

func GetTaskList(ctx *gin.Context) {
	var req taskpb.TaskShowRequest
	// if err := ctx.ShouldBind(&req)
	// no need for bind: need user_id only

	// gin.context has the kv (in jwt)
	user, err := jsonres.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "获取用户信息错误"))
		return
	}

	req = taskpb.TaskShowRequest{
		UserId: user.Id,
	}

	// rpc
	r, err := rpc.TaskList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "TaskShow RPC调用错误"))
		return
	}

	// given struct as JSON into the response body
	// It also sets the Content-Type as "application/json"
	ctx.JSON(http.StatusOK, jsonres.RespSuccess(ctx, r))
}

// TaskCommonResponse
func CreateTask(ctx *gin.Context) {
	var req taskpb.TaskCreateRequest
	// omit blank
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			jsonres.RespErr(ctx, err, "绑定参数错误"))
		return
	}

	// 从 *gin.Context 当中获取到 user_id
	user, err := jsonres.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "获取用户信息错误"))
		return
	}

	// 设置 req 的 userid(也可以直接使用 req 当中的)
	req.UserId = user.Id
	r, err := rpc.TaskCreate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "TaskShow RPC调用错误"))
		return
	}
	ctx.JSON(http.StatusOK, jsonres.RespSuccess(ctx, r))
}

// TaskCommonResponse
func UpdateTask(ctx *gin.Context) {
	var req taskpb.TaskUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			jsonres.RespErr(ctx, err, "绑定参数错误"))
		return
	}

	// unnecessary
	user, err := jsonres.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "获取用户信息错误"))
		return
	}

	if req.UserId != user.Id {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "wrong userid"))
		return
	}

	// rpc
	r, err := rpc.TaskUpdate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "TaskShow RPC调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, jsonres.RespSuccess(ctx, r))
}

// TaskCommonResponse
func DeleteTask(ctx *gin.Context) {
	var req taskpb.TaskDeleteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,
			jsonres.RespErr(ctx, err, "绑定参数错误"))
		return
	}

	// get
	user, err := jsonres.GetUserInfo(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "获取用户信息错误"))
		return
	}

	if req.UserId != user.Id {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "wrong userid"))
		return
	}

	r, err := rpc.TaskDelete(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			jsonres.RespErr(ctx, err, "TaskShow RPC调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, jsonres.RespSuccess(ctx, r))
}
