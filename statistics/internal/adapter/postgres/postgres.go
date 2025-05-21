package postgres

import (
	"context"
	"taskManager/statistics/internal/adapter/postgres/dao"
	"taskManager/statistics/internal/model"
	"taskManager/statistics/pkg/postgre"
)

type Statistics struct {
	DB postgre.Database
}

func New(db postgre.Database) *Statistics {
	return &Statistics{
		DB: db,
	}
}

func (s *Statistics) GetUserStatistics(ctx context.Context) (*model.UserStatistics, error) {
	var daoUser dao.UserStatistics

	err := s.DB.GetDB().First(&daoUser).Error
	if err != nil {
		return nil, err
	}

	returnModel := &model.UserStatistics{
		ID:         daoUser.ID,
		TotalUsers: daoUser.TotalUsers,
	}

	return returnModel, nil
}

func (s *Statistics) GetTaskStatistics(ctx context.Context) (*model.TaskStatistics, error) {
	var daoTask dao.TaskStatistics

	err := s.DB.GetDB().First(&daoTask).Error
	if err != nil {
		return nil, err
	}

	returnModel := &model.TaskStatistics{
		ID:         daoTask.ID,
		TotalTasks: daoTask.TotalTasks,
	}

	return returnModel, nil
}
