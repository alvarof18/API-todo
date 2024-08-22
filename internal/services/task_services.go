package services

import (
	"errors"
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

type TaskServiceImpl struct{}

// All Bussines logic here
var tasks = []models.Task{}

func (s *TaskServiceImpl) GetAllTasks() []models.Task {
	return tasks
}

// Nota coloco el models como puntero para poder devolver un nil porque sino tendria que devolver una structura vacia
func (s *TaskServiceImpl) AddTasks(payload models.Task) (*models.Task, error) {
	newTask := models.Task{
		ID:        uuid.New().String(),
		Title:     payload.Title,
		CreatedAt: time.Now(),
		Completed: false,
	}
	tasks = append(tasks, newTask)
	return &newTask, nil
}

func (s *TaskServiceImpl) FindTaskById(id string) (*models.Task, error) {
	index, err := getIndexTask(id)
	if err != nil {
		return nil, err
	}
	return &tasks[index], nil
}

func (s *TaskServiceImpl) DeleteTask(id string) error {
	index, err := getIndexTask(id)
	if err != nil {
		return err
	}
	tasks = append(tasks[:index], tasks[index+1:]...)
	return nil

}

func (s *TaskServiceImpl) UpdateTask(id string, payload models.UpdateTaskInput) (*models.Task, error) {
	index, err := getIndexTask(id)
	if err != nil {
		return nil, errors.New("task not found to Update")
	}
	updateStruct(&tasks[index], &payload)
	return &tasks[index], nil
}

// No controllers
func getIndexTask(id string) (int, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			return i, nil
		}
	}
	return -1, errors.New("task not found")
}

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
