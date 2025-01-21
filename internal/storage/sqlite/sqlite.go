package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/suvasis-patra/student-api/internal/types"

	"github.com/suvasis-patra/student-api/internal/config"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS student(
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			age INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stm, err := s.Db.Prepare("INSERT INTO student (name,email,age) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stm.Close()
	result, err := stm.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

func (s *Sqlite) FindStudentById(id int64) (types.Student, error) {
	stm, err := s.Db.Prepare("SELECT id,name,email,age FROM student WHERE id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stm.Close()
	var student types.Student
	err = stm.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %v", id)
		}
		return types.Student{}, fmt.Errorf("query error %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stm, err := s.Db.Prepare("SELECT id,name,email,age FROM student")
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	var students []types.Student
	rows, err := stm.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student types.Student
		rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) UpdateStudentDetails(id int64,name string,email string, age int)(int64,error){
	stm,err:=s.Db.Prepare("UPDATE student SET name=?,email=?,age=? WHERE id=?")
	if err!=nil{
		return 0,err
	}
	defer stm.Close()
	res,err:=stm.Exec(name,email,age)
	if err!=nil{
		return 0,err
	}
	rows,err:=res.RowsAffected()
	if err!=nil{
		return 0,err
	}
	return rows,nil
}

func (s *Sqlite) DeleteStudent(id int64) error {
	stm, err := s.Db.Prepare("DELETE FROM student WHERE id=?")
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stm.Close()

	res, err := stm.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no student found with id %d", id)
	}

	return nil
}
