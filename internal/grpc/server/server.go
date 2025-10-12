package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/P04KA/API/gen/stats"
	"github.com/P04KA/API/internal/repository"
	"google.golang.org/grpc"
)

type UserStatsServer struct {
	stats.UnimplementedUserStatsServiceServer
	userRepo repository.UserProvider
}

func NewUserStatsServer(userRepo repository.UserProvider) *UserStatsServer {
	return &UserStatsServer{
		userRepo: userRepo,
	}
}

func (s *UserStatsServer) GetUserStats(ctx context.Context, req *stats.UserStatsRequest) (*stats.UserStatsResponse, error) {

	startTime := time.Now()
	switch req.GetPeriod() {
	case "hour":
		startTime = startTime.Add(-1 * time.Hour)
	case "day":
		startTime = startTime.Add(-24 * time.Hour)
	case "week":
		startTime = startTime.Add(-7 * 24 * time.Hour)
	default:
		startTime = startTime.Add(-24 * time.Hour)
	}

	created, _ := s.userRepo.CountUserCreated(ctx, startTime)
	updated, _ := s.userRepo.CountUserUpdated(ctx, startTime)
	deleted, _ := s.userRepo.CountUserDeleted(ctx, startTime)

	return &stats.UserStatsResponse{
		UsrCreated: created,
		UsrUpdated: updated,
		UsrDeleted: deleted,
		Period:     req.GetPeriod(),
	}, nil
}

func (s *UserStatsServer) Run(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	stats.RegisterUserStatsServiceServer(server, s)

	log.Printf("gRPC сервер запущен на порту %s", port)
	return server.Serve(lis)
}
