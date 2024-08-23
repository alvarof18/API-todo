package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-api/internal/controllers"
	"todo-api/internal/models"
	"todo-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskServices struct {
	mock.Mock
}

// Implementación de los métodos simulados de MockTaskService.
func (m *MockTaskServices) GetAllTasks() []models.Task {
	args := m.Called()
	return args.Get(0).([]models.Task)
}

func (m *MockTaskServices) AddTasks(payload models.Task) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *MockTaskServices) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskServices) FindTaskById(id string) (*models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskServices) UpdateTask(newTask models.Task) error {
	args := m.Called(newTask)
	return args.Error(0)
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router

}

// Prueba para el método GetAllTasks del controlador
func TestGetAllTasks(t *testing.T) {
	taskRepository := new(MockTaskServices)
	taskServices := services.NewTaskService(taskRepository)
	taskController := controllers.NewTaskController(*taskServices)

	//Configurar repuesta simulado

	mockTasks := []models.Task{
		{
			ID:        uuid.New().String(),
			Title:     "Test 1",
			Completed: false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Title:     "Test 2",
			Completed: false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Title:     "Test 2",
			Completed: false,
			CreatedAt: time.Now(),
		},
	}

	taskRepository.On("GetAllTasks").Return(mockTasks)
	router := SetUpRouter()

	router.GET("/tasks", taskController.GetAllTasks)

	// Crear una solicitud simulada
	req, _ := http.NewRequest("GET", "/tasks", nil)

	// Ejecutar la solicitud
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Aserciones
	assert.Equal(t, http.StatusOK, w.Code)

	var responseTasks []models.Task
	json.Unmarshal(w.Body.Bytes(), &responseTasks)

	assert.Len(t, responseTasks, 3)
	taskRepository.AssertExpectations(t)

}

func TestCreateTaskWithMock(t *testing.T) {
	taskRepository := new(MockTaskServices)
	taskServices := services.NewTaskService(taskRepository)
	taskController := controllers.NewTaskController(*taskServices)

	//Configuracion del mock
	taskInput := models.Task{Title: "New Task"}

	mockTask := &models.Task{
		ID:        uuid.New().String(),
		Title:     taskInput.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// Mock the behavior
	taskRepository.On("AddTasks", taskInput).Return(mockTask, nil)

	router := gin.Default()
	router.POST("/task", taskController.AddTasks)

	jsonValue, _ := json.Marshal(taskInput)

	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	taskRepository.AssertExpectations(t)
}

// Test para DeleteTask
func TestDeleteTask(t *testing.T) {
	taskRepository := new(MockTaskServices)
	taskServices := services.NewTaskService(taskRepository)
	taskController := controllers.NewTaskController(*taskServices)
	router := SetUpRouter()

	// Configuración del mock
	taskRepository.On("DeleteTask", "1").Return(nil)

	router.DELETE("/tasks/:id", taskController.DeleteTask)

	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	taskRepository.AssertExpectations(t)
}

// Test para FindTaskById

func TestFindTaskById(t *testing.T) {
	taskRepository := new(MockTaskServices)
	taskServices := services.NewTaskService(taskRepository)
	taskController := controllers.NewTaskController(*taskServices)
	router := SetUpRouter()

	// Configuración del mock
	mockTask := &models.Task{
		ID:        "1",
		Title:     "Task 1",
		CreatedAt: time.Now(),
		Completed: false,
	}
	taskRepository.On("FindTaskById", "1").Return(mockTask, nil)

	router.GET("/tasks/:id", taskController.GetTaskById)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseTask models.Task
	json.Unmarshal(w.Body.Bytes(), &responseTask)

	assert.Equal(t, mockTask.Title, responseTask.Title)

	taskRepository.AssertExpectations(t)
}

// Test para UpdateTask
func TestUpdateTask(t *testing.T) {
	taskRepository := new(MockTaskServices)
	taskServices := services.NewTaskService(taskRepository)
	taskController := controllers.NewTaskController(*taskServices)
	router := SetUpRouter()

	// Configuración del mock
	title := "Updated Task"

	input := models.UpdateTaskInput{
		Title: &title,
	}

	mockTask := &models.Task{
		ID:        "1",
		Title:     *input.Title,
		CreatedAt: time.Now(),
		Completed: false,
	}
	taskRepository.On("UpdateTask", "1", input).Return(mockTask, nil)

	router.PUT("/tasks/:id", taskController.UpdateTask)

	taskJSON, _ := json.Marshal(input)
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseTask models.Task
	json.Unmarshal(w.Body.Bytes(), &responseTask)

	assert.Equal(t, mockTask.Title, responseTask.Title)

	taskRepository.AssertExpectations(t)
}
