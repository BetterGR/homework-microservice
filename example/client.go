package main

import (
	"context"
	"flag"
	"log"
	"time"

	gpb "github.com/BetterGR/homework-microservice/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addr = "localhost:1234"
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := gpb.NewHomeworkServiceClient(conn)

	// Call CreateHomework.
	createHomework(client, "02340311", "Assignment 1", "Design Api and architecture")
	createHomework(client, "02340311", "Assignment 2", "build a boiler plate for the project")

	// Call GetHomework.
	getHomework(client, "CS101")
}

// createHomework is a function that is responsible for making a CreateHomework Request.
func createHomework(client gpb.HomeworkServiceClient, courseID, title, description string) {
	// Create a timeout context for the RPC call.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// Build the request.
	req := &gpb.CreateHomeworkRequest{
		CourseId:    courseID,
		Title:       title,
		Description: description,
	}

	// Make the RPC call.
	res, err := client.CreateHomework(ctx, req)
	if err != nil {
		cancel()
		log.Fatalf("Failed to create homework: %v", err)
	}

	defer cancel()

	// Log the response.
	log.Printf("Homework created successfully for course %s: %v", courseID, res.GetRes())
}

// getHomework is a function that is responsible for making the Get Homework Request.
func getHomework(client gpb.HomeworkServiceClient, courseID string) {
	// Create a timeout context for the RPC call.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// Build the request.
	req := &gpb.GetHomeworkRequest{CourseId: courseID}

	// Make the RPC call.
	res, err := client.GetHomework(ctx, req)
	if err != nil {
		cancel()
		log.Fatalf("Failed to get homework: %v", err)
	}

	defer cancel()

	// Log the response.
	log.Println("Homework List:")

	for _, hw := range res.GetHw() {
		log.Printf("- %s: %s\n", hw.GetTitle(), hw.GetDescription())
	}
}
