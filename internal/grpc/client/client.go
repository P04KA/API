package client

import (
	"context"

	"github.com/P04KA/API/gen/stats"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StatsClient struct {
	client stats.UserStatsServiceClient
	conn   *grpc.ClientConn
}

func New(serverAddr string) (*StatsClient, error) {
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	client := stats.NewUserStatsServiceClient(conn)

	return &StatsClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *StatsClient) GetUserStats(ctx context.Context, period string) (*stats.UserStatsResponse, error) {
	return c.client.GetUserStats(ctx, &stats.UserStatsRequest{
		Period: period,
	})
}

func (c *StatsClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func GetStats(ctx context.Context, period string) (*stats.UserStatsResponse, error) {
	client, err := New("localhost:50052")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.GetUserStats(ctx, period)
}
