package jsonres

import (
	"github.com/gin-gonic/gin"
	"github.com/lukkaal/GRPC-API/pkg/errcode"
)

// standard reply for gateway
// usually no "code" parameter：...int varies from 0-n
func RespSuccess(ctx *gin.Context,
	data interface{}, code ...int) *Response {
	status := errcode.SUCCESS
	if code != nil {
		status = code[0]
	}

	r := &Response{
		Status: status,
		Data:   data,
		Msg:    errcode.GetMsg(status),
	}

	return r
}

func RespErr(ctx *gin.Context, err error,
	data interface{}, code ...int) *Response {
	status := errcode.ERROR
	if code != nil {
		status = code[0]
	}

	r := &Response{
		Status: status,
		Data:   data,
		Msg:    errcode.GetMsg(status),
		Error:  err.Error(),
	}

	return r
}
