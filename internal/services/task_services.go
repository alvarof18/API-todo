package services

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
	"todo-api/internal/models"

	"github.com/google/uuid"
)

type TaskService interface {
	GetAllTasks() []models.Task
	FindTaskById(id string) (*models.Task, error)
	AddTasks(payload models.Task) (*models.Task, error)
	DeleteTask(id string) error
	UpdateTask(id string, payload models.UpdateTaskInput) (*models.Task, error)
}

type TaskServiceImpl struct {
	db *sql.DB
}

// All Bussines logic here

func NewTaskService(db *sql.DB) *TaskServiceImpl {
	return &TaskServiceImpl{db: db}
}

func (s *TaskServiceImpl) GetAllTasks() []models.Task {
	var tasksResult = []models.Task{}

	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return tasksResult
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
			fmt.Printf("Error Getting task %v", err)
			return tasksResult
		}
		tasksResult = append(tasksResult, task)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error Getting tasks %v", err)
		return tasksResult
	}
	return tasksResult
}

// Nota coloco el models como puntero para poder devolver un nil porque sino tendria que devolver una structura vacia
func (s *TaskServiceImpl) AddTasks(payload models.Task) (*models.Task, error) {
	newTask := models.Task{
		ID:        uuid.New().String(),
		Title:     payload.Title,
		CreatedAt: time.Now(),
		Completed: false,
	}
	_, err := s.db.Exec("INSERT INTO tasks VALUES(?,?,?,?)", newTask.ID, newTask.Title, newTask.Completed, newTask.CreatedAt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	return &newTask, nil
}

func (s *TaskServiceImpl) FindTaskById(id string) (*models.Task, error) {

	var task models.Task
	result := s.db.QueryRow("SELECT * FROM tasks where id =?", id)

	if err := result.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %q: no found", id)
		}
		return nil, fmt.Errorf(`error %v`, err)
	}
	return &task, nil
}

func (s *TaskServiceImpl) DeleteTask(id string) error {

	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("task not found to delete")
	}
	fmt.Printf("Rows affected %v", rows)

	return nil
}

func (s *TaskServiceImpl) UpdateTask(id string, payload models.UpdateTaskInput) (*models.Task, error) {

	task, err := s.FindTaskById(id)
	if err != nil {
		return nil, err
	}

	updateStruct(task, &payload)
	_, err = s.db.Exec("UPDATE tasks SET title =?, completed =? WHERE id=?", task.Title, task.Completed, task.ID)
	if err != nil {
		fmt.Printf("Error Update task: %v\n", err)
		return nil, err
	}

	return task, nil
}

// // No controllers
// func getIndexTask(id string) (int, error) {
// 	for i := range tasks {
// 		if tasks[i].ID == id {
// 			return i, nil
// 		}
// 	}
// 	return -1, errors.New("task not found")
// }

// Actualizar informacion usando reflect
func updateStruct(original, input interface{}) {
	// Reflexión sobre los valores originales e input
	origVal := reflect.ValueOf(original).Elem()
	inputVal := reflect.ValueOf(input).Elem()

	for i := 0; i < inputVal.NumField(); i++ {
		// Verificar si el campo es un puntero y no es nil
		field := inputVal.Field(i)
		if !field.IsNil() {
			// Actualiza el campo correspondiente en el struct original
			origField := origVal.FieldByName(inputVal.Type().Field(i).Name)
			origField.Set(reflect.Indirect(field))
		}
	}
}
