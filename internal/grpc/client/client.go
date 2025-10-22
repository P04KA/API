package client

import (
	"context"
	"time"

	"github.com/P04KA/API/pkg/stats"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type statsClient struct {
	client stats.UserStatsServiceClient
	conn   *grpc.ClientConn
}

func New(serverAddr string) (StatsClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &statsClient{
		client: stats.NewUserStatsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *statsClient) GetUserStats(ctx context.Context, period string) (*stats.UserStatsResponse, error) {
	return c.client.GetUserStats(ctx, &stats.UserStatsRequest{
		Period: period,
	})
}

func (c *statsClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
