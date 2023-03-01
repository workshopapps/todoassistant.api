package reminderRepo

import (
	"context"
	"test-va/internals/entity/taskEntity"
	"test-va/internals/entity/vaEntity"
)

type ReminderRepository interface {
	Persist(ctx context.Context, req *taskEntity.CreateTaskReq) error
	PersistAndAssign(ctx context.Context, req *taskEntity.CreateTaskReq) error
	GetPendingTasks(userId string, ctx context.Context) ([]*taskEntity.GetPendingTasksRes, error)
	GetAllUsersPendingTasks() ([]taskEntity.GetPendingTasks, error)
	GetTaskByID(ctx context.Context, taskId string) (*taskEntity.GetTasksByIdRes, error)
	SearchTasks(title *taskEntity.SearchTitleParams, ctx context.Context) ([]*taskEntity.SearchTaskRes, error)
	GetListOfExpiredTasks(ctx context.Context) ([]*taskEntity.GetAllExpiredRes, error)
	GetListOfPendingTasks(ctx context.Context) ([]*taskEntity.GetAllPendingRes, error)
	SetTaskToExpired(id string) error

	GetAllTasks(ctx context.Context, userId string) ([]*taskEntity.GetAllTaskRes, error)
	DeleteTaskByID(ctx context.Context, taskId string) error
	DeleteAllTask(ctx context.Context, userId string) error
	UpdateTaskStatusByID(ctx context.Context, taskId string, req *taskEntity.UpdateTaskStatus) error
	EditTaskById(ctx context.Context, taskId string, req *taskEntity.EditTaskReq) error
	CreateNewTask(req *taskEntity.CreateTaskReq) error

	//VA
	GetAllTaskAssignedToVA(ctx context.Context, vaId string) ([]*vaEntity.VATask, error)
	GetAllTaskForVA(ctx context.Context) ([]*vaEntity.VATaskAll, error)
	GetVADetails(ctx context.Context, userId string) (string, error)
	AssignTaskToVa(ctx context.Context, vaId, taskId string) error

	//Comment
	PersistComment(ctx context.Context, req *taskEntity.CreateCommentReq) error
	GetAllComments(ctx context.Context, taskId string) ([]*taskEntity.GetCommentRes, error)
	GetComments(ctx context.Context) ([]*taskEntity.GetCommentRes, error)
	DeleteCommentByID(ctx context.Context, commentId string) error
}
