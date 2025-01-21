package storage

import "github.com/suvasis-patra/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	FindStudentById(id int64) (types.Student,error)
	GetStudents()([]types.Student,error)
	UpdateStudentDetails(id int64,name string, email string, age int)(int64,error)
	DeleteStudent(id int64)error
}