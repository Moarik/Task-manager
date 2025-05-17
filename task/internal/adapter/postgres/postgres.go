package postgres

import (
	"context"
	"fmt"
	"taskManager/task/internal/adapter/postgres/dao"
	"taskManager/task/internal/model"
	"taskManager/task/pkg/postgre"
)

type Task struct {
	DB postgre.Database
}

func New(db postgre.Database) *Task {
	return &Task{
		DB: db,
	}
}

// CreateUserTask user should be registered before creation
func (t *Task) CreateUserTask(ctx context.Context, task model.Task) (*model.Task, error) {
	insertObject := dao.FromTask(task)

	err := t.DB.GetDB().WithContext(ctx).Create(insertObject).Error
	if err != nil {
		return nil, err
	}

	task.ID = insertObject.ID

	return &task, nil
}

// GetUserTaskByID user should be registered and get only his tasks
func (t *Task) GetUserTaskByID(ctx context.Context, userId, taskId int64) (*model.Task, error) {
	var daoClient dao.Task

	err := t.DB.GetDB().WithContext(ctx).Where("user_id = $1 AND id = $2", userId, taskId).First(&daoClient).Error
	if err != nil {
		return nil, err
	}

	return dao.FromDao(daoClient), nil
}

// DeleteUserTaskByID user should be registered and delete task by his task id
func (t *Task) DeleteUserTaskByID(ctx context.Context, userID, taskID int64) error {
	fmt.Println("USER ID FROM REPO", userID)
	fmt.Println("TASK ID FROM REPO", taskID)

	return t.DB.GetDB().
		WithContext(ctx).
		Unscoped().
		Where("user_id = ? AND id = ?", taskID, userID).
		Delete(&dao.Task{}).Error
}

// UpdateUserTaskByID update users task
func (t *Task) UpdateUserTaskByID(ctx context.Context, task model.Task) (*model.Task, error) {
	var daoTask dao.Task

	if err := t.DB.GetDB().WithContext(ctx).Where("id = $1", task.ID).Delete(&daoTask).Error; err != nil {
		return nil, err
	}

	if task.Title != "" {
		daoTask.Title = task.Title
	}

	if task.Description != "" {
		daoTask.Description = task.Description
	}

	if task.Completed == true {
		daoTask.Completed = true
	}

	err := t.DB.GetDB().WithContext(ctx).Save(&daoTask).Error
	if err != nil {
		return nil, err
	}

	return dao.FromDao(daoTask), nil
}

// GetAllUsersTasks users tasks by his id
func (t *Task) GetAllUsersTasks(ctx context.Context, id int64) ([]*model.Task, error) {
	var daoClient []dao.Task

	err := t.DB.GetDB().WithContext(ctx).Where("user_id = ?", id).Find(&daoClient).Error
	if err != nil {
		return nil, err
	}

	return dao.FromDaos(daoClient), nil
}

// GetTaskByID any task by its id
func (t *Task) GetTaskByID(ctx context.Context, id int64) (*model.Task, error) {
	var daoClient dao.Task

	err := t.DB.GetDB().WithContext(ctx).Where("id = $1", id).First(&daoClient).Error
	if err != nil {
		return nil, err
	}

	return dao.FromDao(daoClient), nil
}
