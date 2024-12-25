package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	gpb "github.com/BetterGR/homework-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	address  = "localhost:1234"
	protocol = "tcp"
)

// homeworkServer is an implementation for the Grpc homework microservice.
type homeworkServer struct {
	// throws unimplemented exception.
	ms.BaseServiceServer
	gpb.UnimplementedHomeworkServiceServer
}

func createHomeworkMicroservice() (*homeworkServer, error) {
	base, err := ms.CreateBaseServiceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base service: %w", err)
	}

	return &homeworkServer{
		BaseServiceServer:                  base,
		UnimplementedHomeworkServiceServer: gpb.UnimplementedHomeworkServiceServer{},
	}, nil
}

// GetHomework handles requests for getting all available homeworks in a course.
func (server *homeworkServer) GetHomework(
	ctx context.Context, req *gpb.GetHomeworkRequest,
) (*gpb.GetHomeworkResponse, error) {
	// Validate the token
	_, err := server.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	homeworks := []*gpb.Homework{
		{Id: "1", Title: "Hw1", Description: "implement bubble sort"},
		{Id: "2", Title: "Hw2", Description: "implement Dijkstra's algorithm"},
	}
	logger.V(0).Info("Fetching homeworks", "courseID", req.GetCourseId(), "homeworkCount", len(homeworks))

	return &gpb.GetHomeworkResponse{Hw: homeworks}, nil
}

// CreateHomework handles requests for adding new homework to a course.
func (server *homeworkServer) CreateHomework(
	ctx context.Context, req *gpb.CreateHomeworkRequest,
) (*gpb.CreateHomeworkResponse, error) {
	logger := klog.FromContext(ctx)
	logger.V(0).Info("New Homework added",
		"courseID", req.GetCourseId(),
		"title", req.GetTitle(),
		"description", req.GetDescription(),
	)

	return &gpb.CreateHomeworkResponse{Res: true}, nil
}

func main() {
	// Initialize the logger.
	klog.InitFlags(nil)
	flag.Parse()

	service, err := createHomeworkMicroservice()
	if err != nil {
		log.Fatalf("Failed to create homework microservice: %v", err)
	}
	// Create a TCP listener.
	lis, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server.
	grpcServer := grpc.NewServer()

	// Register the HomeworkServiceServer with the gRPC server.
	gpb.RegisterHomeworkServiceServer(grpcServer, service)

	klog.Info("gRPC server is running", "address", address)

	if err = grpcServer.Serve(lis); err != nil {
		klog.Error(err, "Failed to serve")
	}
}
