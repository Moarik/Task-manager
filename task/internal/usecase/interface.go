package usecase

import (
	"context"
	"taskManager/task/internal/model"
)

type TaskRepo interface {
	CreateUserTask(ctx context.Context, task model.Task) (*model.Task, error)
	GetUserTaskByID(ctx context.Context, userId, taskId int64) (*model.Task, error)
	DeleteUserTaskByID(ctx context.Context, userId, taskId int64) error
	UpdateUserTaskByID(ctx context.Context, task model.Task) (*model.Task, error)
	GetAllUsersTasks(ctx context.Context, id int64) ([]*model.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*model.Task, error)
}

type TaskEventStorage interface {
	Push(ctx context.Context, user model.Task) error
}
