package taskdb

import (
	"context"

	"github.com/lukkaal/GRPC-API/app/task/internal/repository/taskmodel"
	taskpb "github.com/lukkaal/GRPC-API/idl/task"
	"github.com/lukkaal/GRPC-API/pkg/utils/logger"
	"gorm.io/gorm"
)

type TaskStore struct {
	*gorm.DB // anonimous embeding
}

func NewTaskStore(ctx context.Context) *TaskStore {
	return &TaskStore{
		NewDBClient(ctx),
	}
}

// show task
func (taskstore *TaskStore) TashShowList(userid int64) (
	tasks []*taskmodel.Task, err error) {
	err = taskstore.Model(&taskmodel.Task{}).
		Where("user_id=?", userid).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return
}

// create task(primary key: taskid -> auto generated)
func (taskstore *TaskStore) CreateTask(
	req *taskpb.TaskCreateRequest) (err error) {
	// no need for checking
	newTask := &taskmodel.Task{
		UserID:    req.UserId,
		Title:     req.Title,
		Content:   req.Content,
		Status:    int(req.Status),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	if err = taskstore.Model(&taskmodel.Task{}).
		Create(newTask).Error; err != nil {
		logger.GinloggerObj.Error("Create Error:" + err.Error())
		return err
	}

	return nil
}

// delete task
func (taskstore *TaskStore) DeleteTaskById(
	taskId, userId int64) (err error) {
	err = taskstore.Model(&taskmodel.Task{}).
		Where("task_id = ? AND user_id = ?",
			taskId, userId).Error

	return
}

// update by taskid
func (taskstore *TaskStore) UpdateTask(
	req *taskpb.TaskUpdateRequest) (err error) {
	taskUpdateMap := make(map[string]interface{})

	taskUpdateMap["title"] = req.Title
	taskUpdateMap["content"] = req.Content
	taskUpdateMap["status"] = int(req.Status)
	taskUpdateMap["start_time"] = req.StartTime
	taskUpdateMap["end_time"] = req.EndTime
	taskUpdateMap["user_id"] = req.UserId

	err = taskstore.Model(&taskmodel.Task{}).
		Where("task_id = ? AND user_id = ?",
			req.TaskId, req.UserId).Updates(&taskUpdateMap).Error

	return
}
