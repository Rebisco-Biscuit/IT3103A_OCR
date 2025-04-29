package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection pool
var DB *sql.DB

// Initialize the database connection
func InitDB() {
	// Replace with your actual PostgreSQL connection string
	connStr := "user=postgres password=mysecretpassword dbname=student_service sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verify the connection is successful
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connected successfully")
}

// GetStudentByID fetches a student by their ID from the database
func GetStudentByID(id string) (*Student, error) {
	var student Student
	query := "SELECT id, name, birthdate FROM students WHERE id = $1"
	err := DB.QueryRow(query, id).Scan(&student.ID, &student.Name, &student.Birthdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no student found with id: %s", id)
		}
		return nil, fmt.Errorf("error querying student by ID: %v", err)
	}
	return &student, nil
}

// CreateStudent adds a new student to the database
func CreateStudent(name string, birthdate string) (*Student, error) {
	var student Student
	query := "INSERT INTO students (name, birthdate) VALUES ($1, $2) RETURNING id"
	err := DB.QueryRow(query, name, birthdate).Scan(&student.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating student: %v", err)
	}
	student.Name = name
	student.Birthdate = birthdate
	return &student, nil
}

// Student struct represents a student entity in the system
type Student struct {
	ID        string
	Name      string
	Birthdate string
}
