package rpc

import (
	"context"
	"errors"

	taskpb "github.com/lukkaal/GRPC-API/idl/task"
	"github.com/lukkaal/GRPC-API/pkg/errcode"
)

// seal the grpc func into gin.Handler(1st step)
func TaskCreate(ctx context.Context,
	req *taskpb.TaskCreateRequest) (
	resp *taskpb.TaskCommonResponse, err error) {

	resp, err = TaskClient.TaskCreate(ctx, req)
	if err != nil {
		return
	}

	if resp.Code != errcode.SUCCESS {
		err = errors.New(resp.Msg)
	}

	return
}

func TaskUpdate(ctx context.Context,
	req *taskpb.TaskUpdateRequest) (
	resp *taskpb.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskUpdate(ctx, req)
	if err != nil {
		return
	}

	if resp.Code != errcode.SUCCESS {
		err = errors.New(resp.Msg)
		return
	}

	return
}

func TaskDelete(ctx context.Context,
	req *taskpb.TaskDeleteRequest) (
	resp *taskpb.TaskCommonResponse, err error) {
	resp, err = TaskClient.TaskDelete(ctx, req)

	if err != nil {
		return
	}

	if resp.Code != errcode.SUCCESS {
		err = errors.New(resp.Msg)
		return
	}

	return
}

func TaskList(ctx context.Context,
	req *taskpb.TaskShowRequest) (
	resp *taskpb.TasksDetailResponse, err error) {
	resp, err = TaskClient.TaskShow(ctx, req)

	if err != nil {
		return
	}

	if resp.Code != errcode.SUCCESS {
		err = errors.New("fetch fail")
		return
	}

	return
}
