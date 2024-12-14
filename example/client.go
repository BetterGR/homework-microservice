package main

import (
	"context"
	"flag"
	"log"
	"time"

	gpb "github.com/BetterGR/homework-microservice/homework_protos" // Replace with your package path
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:1234", "the address to connect to") // Set port to 1234
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := gpb.NewHomeworkServiceClient(conn)

	// Call CreateHomework
	createHomework(client, "02340311", "Assignment 1", "Design Api and architecture")
	createHomework(client, "02340311", "Assignment 2", "build a boiler plate for the project")

	// Call GetHomework
	getHomework(client, "CS101")
}

// createHomework is a function that is responsible for making a CreateHomework Request
func createHomework(client gpb.HomeworkServiceClient, courseId, title, description string) {
	// Create a timeout context for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Build the request
	req := &gpb.CreateHomeworkRequest{
		CourseId:    courseId,
		Title:       title,
		Description: description,
	}

	// Make the RPC call
	res, err := client.CreateHomework(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create homework: %v", err)
	}

	// Log the response
	log.Printf("Homework created successfully for course %s: %v", courseId, res.Res)
}

// getHomework is a function that is responsible for making the Get Homework Request
func getHomework(client gpb.HomeworkServiceClient, courseId string) {
	// Create a timeout context for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Build the request
	req := &gpb.GetHomeworkRequest{CourseId: courseId}

	// Make the RPC call
	res, err := client.GetHomework(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get homework: %v", err)
	}

	// Log the response
	log.Println("Homework List:")
	for _, hw := range res.Hw {
		log.Printf("- %s: %s\n", hw.Title, hw.Description)
	}
}
