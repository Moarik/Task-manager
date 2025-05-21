package usecase

import (
	"context"
	"taskManager/statistics/internal/model"
)

type StatisticsRepo interface {
	GetUserStatistics(ctx context.Context) (*model.UserStatistics, error)
	GetTaskStatistics(ctx context.Context) (*model.TaskStatistics, error)
}
