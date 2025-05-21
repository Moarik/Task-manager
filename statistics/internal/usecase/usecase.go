package usecase

import (
	"context"
	"taskManager/statistics/internal/model"
)

type Statistics struct {
	Repo StatisticsRepo
}

func New(repo StatisticsRepo) *Statistics {
	return &Statistics{
		Repo: repo,
	}
}

func (s *Statistics) GetUserStatistics(ctx context.Context) (*model.UserStatistics, error) {
	statistics, err := s.Repo.GetUserStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func (s *Statistics) GetTaskStatistics(ctx context.Context) (*model.TaskStatistics, error) {
	statistics, err := s.Repo.GetTaskStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}
