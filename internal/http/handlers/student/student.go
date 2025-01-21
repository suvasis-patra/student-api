package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/suvasis-patra/student-api/internal/storage"
	"github.com/suvasis-patra/student-api/internal/types"
	"github.com/suvasis-patra/student-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		// general error
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		// request validation
		if validationErr := validator.New().Struct(student); validationErr != nil {
			// type casting
			err := validationErr.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err))
			return
		}
		// create the student in database
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"success": lastId})
	}
}

func GetStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		studentId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.FindStudentById(studentId)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetAllStudents(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		students,err:=storage.GetStudents()
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		response.WriteJson(w,http.StatusOK,students)		
	}
}

func UpdateStudentDetailsById(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		id:=r.PathValue("id")
		var student types.Student
		err:=json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		if err:=validator.New().Struct(student); err!=nil{
			validationErr:=err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError((validationErr)))
			return
		} 
		intId,err:=strconv.ParseInt(id,10,64)
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		studentId,err:=storage.UpdateStudentDetails(intId,student.Name,student.Email,student.Age)
		if err!=nil{
			response.WriteJson(w,http.StatusNotFound,response.GeneralError(fmt.Errorf("no student found with id %v",intId)))
			return
		}
		response.WriteJson(w,http.StatusOK,studentId)
	}
}

func DeleteStudentById(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		id:=r.PathValue("id")
		intId,err:=strconv.ParseInt(id,10,64)
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		delError:=storage.DeleteStudent(intId)
		if delError!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(delError))
			return
		}
		response.WriteJson(w,http.StatusOK,map[string]string{"status": "OK"})
	}
}
