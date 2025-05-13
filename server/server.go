package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	hpb "github.com/BetterGR/homework-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	// define address.
	connectionProtocol = "tcp"
	// Debugging logs.
	logLevelDebug = 5
)

type HomeworkServer struct {
	ms.BaseServiceServer
	db *Database
	// throws unimplemented error
	hpb.UnimplementedHomeworkServiceServer
}

// initHomeworkMicroserviceServer initializes the StudentsServer.
func initHomeworkMicroserviceServer() (*HomeworkServer, error) {
	base, err := ms.CreateBaseServiceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base service: %w", err)
	}

	database, err := InitializeDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &HomeworkServer{
		BaseServiceServer:                  base,
		db:                                 database,
		UnimplementedHomeworkServiceServer: hpb.UnimplementedHomeworkServiceServer{},
	}, nil
}

// CreateHomework creates a new homework.
func (s *HomeworkServer) CreateHomework(ctx context.Context,
	req *hpb.CreateHomeworkRequest,
) (*hpb.CreateHomeworkResponse, error) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received CreateHomework request", "courseId",
		req.GetHomework().GetCourseId(), "title", req.GetHomework().GetTitle())

	homework := req.GetHomework()
	if homework == nil {
		return nil, status.Errorf(codes.InvalidArgument, "homework is nil")
	}

	// insert the homework into the database.

	if err := s.db.AddHomework(ctx, homework); err != nil {
		logger.Error(err, "failed to insert homework")
		return nil, status.Errorf(codes.Internal, "failed to insert homework: %v", err)
	}

	logger.V(logLevelDebug).Info("Successfully created homework", "id", req.GetHomework().GetId())

	return &hpb.CreateHomeworkResponse{Hw: req.GetHomework()}, nil
}

// GetHomework retrieves a homework by ID.
func (s *HomeworkServer) GetHomework(ctx context.Context,
	req *hpb.GetHomeworkRequest,
) (*hpb.GetHomeworkResponse, error) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received GetHomework request", "id", req.GetId())

	// get the homework from the database.
	homework, err := s.db.GetHomework(ctx, req.GetId())
	if err != nil {
		logger.Error(err, "failed to get homework", "id", req.GetId())
		return nil, status.Errorf(codes.NotFound, "failed to get homework: %v", err)
	}

	logger.V(logLevelDebug).Info("Successfully fetched homework", "id", req.GetId())

	return &hpb.GetHomeworkResponse{Hw: homework}, nil
}

// UpdateHomework updates an existing homework.
func (s *HomeworkServer) UpdateHomework(ctx context.Context,
	req *hpb.UpdateHomeworkRequest,
) (*hpb.UpdateHomeworkResponse, error) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received UpdateHomework request", "id", req.GetHomework().GetId())

	homework := req.GetHomework()
	if homework == nil {
		return nil, status.Errorf(codes.InvalidArgument, "homework is nil")
	}

	// update the homework in the database.
	if err := s.db.UpdateHomework(ctx, homework); err != nil {
		logger.Error(err, "failed to update homework", "id", req.GetHomework().GetId())
		return nil, status.Errorf(codes.Internal, "failed to update homework: %v", err)
	}

	logger.V(logLevelDebug).Info("Successfully updated homework", "id", req.GetHomework().GetId())

	return &hpb.UpdateHomeworkResponse{Hw: homework}, nil
}

// DeleteHomework deletes a homework by ID.
func (s *HomeworkServer) DeleteHomework(ctx context.Context,
	req *hpb.DeleteHomeworkRequest,
) (*hpb.DeleteHomeworkResponse, error) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received DeleteHomework request", "id", req.GetId())

	// delete the homework from the database.
	if err := s.db.DeleteHomework(ctx, req.GetId()); err != nil {
		logger.Error(err, "failed to delete homework", "id", req.GetId())
		return nil, status.Errorf(codes.NotFound, "failed to delete homework: %v", err)
	}

	logger.V(logLevelDebug).Info("Successfully deleted homework", "id", req.GetId())

	return &hpb.DeleteHomeworkResponse{}, nil
}

// main StudentsServer function.
func main() {
	// init klog
	klog.InitFlags(nil)
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		klog.Fatalf("Error loading .env file")
	}

	// init the StudentsServer
	server, err := initHomeworkMicroserviceServer()
	if err != nil {
		klog.Fatalf("Failed to init HomeworkServer: %v", err)
	}

	// create a listener on port 'address'
	address := os.Getenv("GRPC_PORT")

	lis, err := net.Listen(connectionProtocol, address)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	klog.Info("Starting Homework on port: ", address)
	// create a grpc HomeworkServer
	grpcServer := grpc.NewServer()
	hpb.RegisterHomeworkServiceServer(grpcServer, server)

	// serve the grpc StudentsServer
	if err := grpcServer.Serve(lis); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
