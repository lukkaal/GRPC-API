package service

import (
	"context"
	"sync"

	"github.com/lukkaal/GRPC-API/app/task/internal/repository/taskdb"
	taskpb "github.com/lukkaal/GRPC-API/idl/task"
	"github.com/lukkaal/GRPC-API/pkg/errcode"
)

// must be embedded to have forward
// compatible implementations
type TaskSrv struct {
	taskpb.UnimplementedTaskServiceServer
}

var TaskSrvIns *TaskSrv
var TaskSrvOnce sync.Once

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvIns = &TaskSrv{}
	})

	return TaskSrvIns
}

// taskcreate
func (*TaskSrv) TaskCreate(ctx context.Context,
	req *taskpb.TaskCreateRequest) (
	resp *taskpb.TaskCommonResponse, err error) {
	resp.Code = errcode.SUCCESS
	err = taskdb.NewTaskStore(ctx).CreateTask(req)
	if err != nil {
		resp.Code = errcode.ERROR
		resp.Data = err.Error()
	}

	resp.Msg = errcode.GetMsg(int(resp.Code))
	return
}

// show task by userid
func (*TaskSrv) TaskShow(ctx context.Context,
	req *taskpb.TaskShowRequest) (
	resp *taskpb.TasksDetailResponse, err error) {
	resp = new(taskpb.TasksDetailResponse)
	resp.Code = errcode.SUCCESS

	tasks, err := taskdb.NewTaskStore(ctx).TashShowList(req.UserId)
	if err != nil {
		resp.Code = errcode.ERROR
		return
	}

	// append to response
	for idx, _ := range tasks {
		resp.TaskDetail = append(resp.TaskDetail,
			&taskpb.TaskModel{
				TaskId:    tasks[idx].TaskID,
				Status:    int64(tasks[idx].Status),
				UserId:    tasks[idx].UserID,
				Content:   tasks[idx].Content,
				Title:     tasks[idx].Title,
				StartTime: tasks[idx].StartTime,
				EndTime:   tasks[idx].EndTime,
			})
	}

	return
}

// update task by taskid
func (*TaskSrv) TaskUpdate(ctx context.Context,
	req *taskpb.TaskUpdateRequest) (
	resp *taskpb.TaskCommonResponse, err error) {
	resp.Code = errcode.SUCCESS
	err = taskdb.NewTaskStore(ctx).UpdateTask(req)
	if err != nil {
		resp.Code = errcode.ERROR
		resp.Data = err.Error()
	}
	resp.Msg = errcode.GetMsg(
		int(resp.Code))

	return
}

// delete task by taskid and userid
func (*TaskSrv) TaskDelete(ctx context.Context,
	req *taskpb.TaskDeleteRequest) (
	resp *taskpb.TaskCommonResponse, err error) {

	resp.Code = errcode.SUCCESS

	err = taskdb.NewTaskStore(ctx).
		DeleteTaskById(req.TaskId, req.UserId)
	if err != nil {
		resp.Code = errcode.ERROR
		resp.Data = err.Error()
	}

	resp.Msg = errcode.GetMsg(
		int(resp.Code))

	return
}
