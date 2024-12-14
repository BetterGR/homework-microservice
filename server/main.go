package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	gpb "github.com/BetterGR/homework-microservice/homework_protos"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	address  = "localhost:1234"
	protocol = "tcp"
)

// homeworkServer is an implementation for the Grpc homework microservice
type homeworkServer struct {
	//throws unimplemented exception
	gpb.UnimplementedHomeworkServiceServer
	logger logr.Logger
}

// newHomeworkServer is a constructor for the homeworkServer
func newHomeworkServer(logger logr.Logger) *homeworkServer {
	return &homeworkServer{logger: logger}
}

// setupLogger initializes the zap-based logr logger
func setupLogger() logr.Logger {
	zapLog, err := zap.NewDevelopment() // Use zap.NewProduction() for production environments
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	return zapr.NewLogger(zapLog)
}

// GetHomework function handles the request for getting all the homeworks that are available in a certain course
func (s *homeworkServer) GetHomework(ctx context.Context, req *gpb.GetHomeworkRequest) (*gpb.GetHomeworkResponse, error) {
	// Example: Hardcoded homework list
	homeworks := []*gpb.Homework{
		{Id: "1", Title: "Hw1", Description: "implement bubble sort"},
		{Id: "2", Title: "Hw2", Description: "implement Dijkstra's algorithm"},
	}
	s.logger.Info("Fetching homeworks", "courseID", req.CourseId, "homeworkCount", len(homeworks))
	return &gpb.GetHomeworkResponse{Hw: homeworks}, nil
}

// CreateHomework function handles the requests for adding new home work to a certain course
func (s *homeworkServer) CreateHomework(ctx context.Context, req *gpb.CreateHomeworkRequest) (*gpb.CreateHomeworkResponse, error) {
	s.logger.Info("New Homework added",
		"courseID", req.CourseId,
		"title", req.Title,
		"description", req.Description,
	)

	return &gpb.CreateHomeworkResponse{Res: true}, nil
}

func main() {
	// Initialize the logger
	logger := setupLogger()
	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Errorf("%v", r), "Application crashed")
			os.Exit(1)
		}
	}()

	// Create a TCP listener
	lis, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	server := newHomeworkServer(logger)

	// Register the HomeworkServiceServer with the gRPC server
	gpb.RegisterHomeworkServiceServer(grpcServer, server)

	logger.Info("gRPC server is running", "address", address)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error(err, "Failed to serve")
	}
}
