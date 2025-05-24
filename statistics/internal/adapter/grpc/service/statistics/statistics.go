package statistics

import (
	"context"
	"taskManager/proto/gen/statistics"
)

type Client struct {
	statistics.UnimplementedStatisticsServiceServer

	uc StatisticsUsecase
}

func New(uc StatisticsUsecase) *Client {
	return &Client{
		uc: uc,
	}
}

func (c *Client) GetUserStatistics(ctx context.Context, req *statistics.Empty) (*statistics.GetUserStatisticsResponse, error) {
	returnTask, err := c.uc.GetUserStatistics(ctx)
	if err != nil {
		return nil, err
	}

	sendTask := &statistics.GetUserStatisticsResponse{
		Id:         returnTask.ID,
		TotalUsers: int32(returnTask.TotalUsers),
	}

	return sendTask, err
}

func (c *Client) GetTaskStatistics(ctx context.Context, req *statistics.Empty) (*statistics.GetTaskStatisticsResponse, error) {
	returnTask, err := c.uc.GetTaskStatistics(ctx)
	if err != nil {
		return nil, err
	}

	sendTask := &statistics.GetTaskStatisticsResponse{
		Id:         returnTask.ID,
		TotalTasks: int32(returnTask.TotalTasks),
	}

	return sendTask, err
}
