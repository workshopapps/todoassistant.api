package reminderRepo

import (
	"test-va/internals/entity/reminderEntity"
	"test-va/internals/entity/taskEntity"
)

type ReminderRepository interface {
	// Needed
	SetTaskToExpired(id string) error
	CreateNewTask(req *taskEntity.CreateTaskReq) error
	GetAllUsersPendingTasks() ([]reminderEntity.GetPendingTasks, error)
}
