package postgres

import (
	"context"
	"errors"
	"gorm.io/gorm"
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

	err := s.DB.GetDB().WithContext(ctx).Unscoped().First(&daoUser).Error
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

	err := s.DB.GetDB().WithContext(ctx).First(&daoTask).Error
	if err != nil {
		return nil, err
	}

	returnModel := &model.TaskStatistics{
		ID:         daoTask.ID,
		TotalTasks: daoTask.TotalTasks,
	}

	return returnModel, nil
}

func (s *Statistics) CreateUserStatistics(ctx context.Context, client model.UserNats) error {
	if client.ID == 0 {
		return errors.New("invalid user ID: 0")
	}

	var stat dao.UserStatistics
	db := s.DB.GetDB().WithContext(ctx)

	err := db.First(&stat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stat = dao.UserStatistics{
			TotalUsers: 1,
		}
		return db.Create(&stat).Error
	} else if err != nil {
		return err
	}

	stat.TotalUsers += 1
	return db.Save(&stat).Error
}

func (s *Statistics) CreateTaskStatistics(ctx context.Context, client model.TaskNats) error {
	if client.ID == 0 {
		return errors.New("invalid task ID: 0")
	}

	var stat dao.TaskStatistics
	db := s.DB.GetDB().WithContext(ctx)

	err := db.First(&stat).Error
	if err != nil {
		return err
	}

	stat.TotalTasks += 1

	return db.Save(&stat).Error
}
