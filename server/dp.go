package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	hpb "github.com/BetterGR/homework-microservice/protos"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"k8s.io/klog/v2"
)

// Database represents the database connection.
type Database struct {
	db *bun.DB
}

// InitializeDatabase ensures that the database exists and initializes the schema.
func InitializeDatabase() (*Database, error) {
	createDatabaseIfNotExists()

	database, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	if err := database.createSchemaIfNotExists(context.Background()); err != nil {
		klog.Fatalf("Failed to create schema: %v", err)
	}

	return database, nil
}

// createDatabaseIfNotExists ensures the database exists.
func createDatabaseIfNotExists() {
	dsn := os.Getenv("DSN")
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))

	sqldb := sql.OpenDB(connector)
	defer sqldb.Close()

	ctx := context.Background()
	dbName := os.Getenv("DP_NAME")
	query := "SELECT 1 FROM pg_database WHERE datname = $1;"

	var exists int

	err := sqldb.QueryRowContext(ctx, query, dbName).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		klog.Fatalf("Failed to check db existence: %v", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		if _, err = sqldb.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s;", dbName)); err != nil {
			klog.Fatalf("Failed to create database: %v", err)
		}

		klog.Infof("Database %s created successfully.", dbName)
	} else {
		klog.Infof("Database %s already exists.", dbName)
	}
}

// ConnectDB connects to the database.
func ConnectDB() (*Database, error) {
	dsn := os.Getenv("DSN")
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)
	database := bun.NewDB(sqldb, pgdialect.New())

	// Test the connection.
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	klog.Info("Connected to PostgreSQL database.")

	return &Database{db: database}, nil
}

// createSchemaIfNotExists creates the database schema if it doesn't exist.
func (d *Database) createSchemaIfNotExists(ctx context.Context) error {
	models := []interface{}{
		(*Homework)(nil),
	}

	for _, model := range models {
		if _, err := d.db.NewCreateTable().IfNotExists().Model(model).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	klog.Info("Database schema initialized.")

	return nil
}

type Homework struct {
	UniqueID    string            `bun:",pk,default:gen_random_uuid()"`
	ID          string            `bun:"id,unique,notnull"`
	CourseID    string            `bun:"course_id,notnull"`
	Title       string            `bun:"title,notnull"`
	Description string            `bun:"description,notnull"`
	Files       []*hpb.File       `bun:"files,array"`
	Workflow    string            `bun:"workflow,notnull"`
	DueDate     string            `bun:"due_date,notnull"`
	Submissions []*hpb.Submission `bun:"submissions,array"`
}

// AddHomework adds a homework to the database.
func (d *Database) AddHomework(ctx context.Context, homework *hpb.Homework) error {
	if _, err := d.db.NewInsert().Model(&Homework{
		ID:          homework.GetId(),
		CourseID:    homework.GetCourseId(),
		Title:       homework.GetTitle(),
		Description: homework.GetDescription(),
		Files:       homework.GetFiles(),
		Workflow:    homework.GetWorkflow(),
		DueDate:     homework.GetDueDate(),
		Submissions: homework.GetSubmissions(),
	}).Exec(ctx); err != nil {
		return fmt.Errorf("failed to insert homework: %w", err)
	}

	klog.Info("Homework added successfully.")

	return nil
}

// GetHomework retrieves a homework by ID from the database.
func (d *Database) GetHomework(ctx context.Context, id string) (*hpb.Homework, error) {
	homework := new(Homework)

	if err := d.db.NewSelect().Model(homework).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to get homework: %w", err)
	}

	return &hpb.Homework{
		Id:          homework.ID,
		CourseId:    homework.CourseID,
		Title:       homework.Title,
		Description: homework.Description,
		Files:       homework.Files,
		Workflow:    homework.Workflow,
		DueDate:     homework.DueDate,
		Submissions: homework.Submissions,
	}, nil
}

// UpdateHomework updates an existing homework in the database.
func (d *Database) UpdateHomework(ctx context.Context, homework *hpb.Homework) error {
	_, err := d.db.NewUpdate().Model(&Homework{
		ID:          homework.GetId(),
		CourseID:    homework.GetCourseId(),
		Title:       homework.GetTitle(),
		Description: homework.GetDescription(),
		Files:       homework.GetFiles(),
		Workflow:    homework.GetWorkflow(),
		DueDate:     homework.GetDueDate(),
		Submissions: homework.GetSubmissions(),
	}).Where("id = ?", homework.GetId()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update homework: %w", err)
	}

	klog.Info("Homework updated successfully.")

	return nil
}

// DeleteHomework removes a homework from the database.
func (d *Database) DeleteHomework(ctx context.Context, id string) error {
	_, err := d.db.NewDelete().Model((*Homework)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete homework: %w", err)
	}

	klog.Info("Homework deleted successfully.")

	return nil
}
