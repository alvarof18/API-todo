package services

import (
	"reflect"
	"time"
	"todo-api/internal/models"
	"todo-api/internal/repositories"

	"github.com/google/uuid"
)

type TaskService struct {
	repo repositories.TaskRepository
}

// All Bussines logic here
func NewTaskService(repo repositories.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAllTasks() []models.Task {
	return s.repo.GetAllTasks()
}

// Nota coloco el models como puntero para poder devolver un nil porque sino tendria que devolver una structura vacia
func (s *TaskService) AddTasks(payload models.Task) (*models.Task, error) {
	newTask := models.Task{
		ID:        uuid.New().String(),
		Title:     payload.Title,
		CreatedAt: time.Now(),
		Completed: false,
	}
	err := s.repo.AddTasks(newTask)
	return &newTask, err
}

func (s *TaskService) FindTaskById(id string) (*models.Task, error) {
	return s.repo.FindTaskById(id)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.repo.DeleteTask(id)
}

func (s *TaskService) UpdateTask(id string, payload models.UpdateTaskInput) (*models.Task, error) {
	task, err := s.repo.FindTaskById(id)
	if err != nil {
		return nil, err
	}

	updateStruct(task, &payload)
	return task, s.repo.UpdateTask(*task)
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
