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

func (m *MockTaskServices) AddTasks(payload models.Task) (*models.Task, error) {
	args := m.Called(payload)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskServices) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskServices) FindTaskById(id string) (*models.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskServices) UpdateTask(id string, payload models.UpdateTaskInput) (*models.Task, error) {
	args := m.Called(id, payload)
	return args.Get(0).(*models.Task), args.Error(1)
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router

}

// Prueba para el método GetAllTasks del controlador
func TestGetAllTasks(t *testing.T) {
	mockServices := new(MockTaskServices)
	taskController := controllers.NewTaskController(mockServices)

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

	mockServices.On("GetAllTasks").Return(mockTasks)
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
	mockServices.AssertExpectations(t)

}

func TestCreateTaskWithMock(t *testing.T) {
	mockService := new(MockTaskServices)
	controller := controllers.NewTaskController(mockService)

	//Configuracion del mock
	taskInput := models.Task{Title: "New Task"}

	mockTask := &models.Task{
		ID:        uuid.New().String(),
		Title:     taskInput.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// Mock the behavior
	mockService.On("AddTasks", taskInput).Return(mockTask, nil)

	router := gin.Default()
	router.POST("/task", controller.AddTasks)

	jsonValue, _ := json.Marshal(taskInput)

	req, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	mockService.AssertExpectations(t)
}

// Test para DeleteTask
func TestDeleteTask(t *testing.T) {
	mockService := new(MockTaskServices)
	router := SetUpRouter()
	controllers := controllers.NewTaskController(mockService)

	// Configuración del mock
	mockService.On("DeleteTask", "1").Return(nil)

	router.DELETE("/tasks/:id", controllers.DeleteTask)

	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockService.AssertExpectations(t)
}

// Test para FindTaskById

func TestFindTaskById(t *testing.T) {
	mockService := new(MockTaskServices)
	controllers := controllers.NewTaskController(mockService)
	router := SetUpRouter()

	// Configuración del mock
	mockTask := &models.Task{
		ID:        "1",
		Title:     "Task 1",
		CreatedAt: time.Now(),
		Completed: false,
	}
	mockService.On("FindTaskById", "1").Return(mockTask, nil)

	router.GET("/tasks/:id", controllers.GetTaskById)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseTask models.Task
	json.Unmarshal(w.Body.Bytes(), &responseTask)

	assert.Equal(t, mockTask.Title, responseTask.Title)

	mockService.AssertExpectations(t)
}

// Test para UpdateTask
func TestUpdateTask(t *testing.T) {
	mockService := new(MockTaskServices)
	controller := controllers.NewTaskController(mockService)
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
	mockService.On("UpdateTask", "1", input).Return(mockTask, nil)

	router.PUT("/tasks/:id", controller.UpdateTask)

	taskJSON, _ := json.Marshal(input)
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseTask models.Task
	json.Unmarshal(w.Body.Bytes(), &responseTask)

	assert.Equal(t, mockTask.Title, responseTask.Title)

	mockService.AssertExpectations(t)
}
