package resolver

import (
	"context"
	"student-service/database"
	"student-service/graphql/model"
)

// Resolver holds the database connection and resolvers for GraphQL queries
type Resolver struct {
	DB *database.DB
}

// GetStudent resolver fetches a student by their ID
func (r *Resolver) GetStudent(ctx context.Context, id string) (*model.Student, error) {
	student, err := r.DB.GetStudentByID(id)
	if err != nil {
		return nil, err
	}
	return student, nil
}

// ListStudents resolver returns a list of students
func (r *Resolver) ListStudents(ctx context.Context) ([]*model.Student, error) {
	students, err := r.DB.GetAllStudents()
	if err != nil {
		return nil, err
	}
	return students, nil
}

// CreateStudent resolver creates a new student
func (r *Resolver) CreateStudent(ctx context.Context, name, birthdate, gender, location, bio, interests, workHistory, educationHistory string) (*model.Student, error) {
	student, err := r.DB.CreateStudent(name, birthdate, gender, location, bio, interests, workHistory, educationHistory)
	if err != nil {
		return nil, err
	}
	return student, nil
}

// UpdateStudent resolver updates a student's details
func (r *Resolver) UpdateStudent(ctx context.Context, id, name, birthdate, gender, location, bio, interests, workHistory, educationHistory string) (*model.Student, error) {
	student, err := r.DB.UpdateStudent(id, name, birthdate, gender, location, bio, interests, workHistory, educationHistory)
	if err != nil {
		return nil, err
	}
	return student, nil
}

// DeleteStudent resolver deletes a student
func (r *Resolver) DeleteStudent(ctx context.Context, id string) (bool, error) {
	err := r.DB.DeleteStudent(id)
	if err != nil {
		return false, err
	}
	return true, nil
}
