package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"flag"

	gpb "github.com/BetterGR/homework-microservice/homework_protos"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

const (
	address  = "localhost:1234"
	protocol = "tcp"
)

// homeworkServer is an implementation for the Grpc homework microservice.
type homeworkServer struct {
	//throws unimplemented exception.
	gpb.UnimplementedHomeworkServiceServer
}

// setupKLogger initializes Klog logger.
func setupKLogger() {
	klog.InitFlags(nil)
	flag.Parse()
}

// GetHomework function handles the request for getting all the homeworks that are available in a certain course.
func (s *homeworkServer) GetHomework(ctx context.Context, req *gpb.GetHomeworkRequest) (*gpb.GetHomeworkResponse, error) {
	// Example: Hardcoded homework list.
	logger := klog.FromContext(ctx)
	homeworks := []*gpb.Homework{
		{Id: "1", Title: "Hw1", Description: "implement bubble sort"},
		{Id: "2", Title: "Hw2", Description: "implement Dijkstra's algorithm"},
	}
	logger.V(0).Info("Fetching homeworks", "courseID", req.CourseId, "homeworkCount", len(homeworks))
	return &gpb.GetHomeworkResponse{Hw: homeworks}, nil
}

// CreateHomework function handles the requests for adding new home work to a certain course.
func (s *homeworkServer) CreateHomework(ctx context.Context, req *gpb.CreateHomeworkRequest) (*gpb.CreateHomeworkResponse, error) {
	logger := klog.FromContext(ctx)
	logger.V(0).Info("New Homework added",
		"courseID", req.CourseId,
		"title", req.Title,
		"description", req.Description,
	)

	return &gpb.CreateHomeworkResponse{Res: true}, nil
}

func main() {
	// Initialize the logger.
	setupKLogger()

	defer func() {
		if r := recover(); r != nil {
			klog.Error(fmt.Errorf("%v", r), "Application crashed")
			os.Exit(1)
		}
	}()

	// Create a TCP listener.
	lis, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server.
	grpcServer := grpc.NewServer()

	// Register the HomeworkServiceServer with the gRPC server.
	gpb.RegisterHomeworkServiceServer(grpcServer, &homeworkServer{})

	klog.Info("gRPC server is running", "address", address)
	if err := grpcServer.Serve(lis); err != nil {
		klog.Error(err, "Failed to serve")
	}
}
