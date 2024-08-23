// Maneja las peticiones HTTP

package controllers

import (
	"net/http"
	"todo-api/internal/models"
	"todo-api/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService services.TaskService
}

// Simular una inyeccion de dependencia
func NewTaskController(service services.TaskService) *TaskController {
	return &TaskController{taskService: service}
}

func (ctrl *TaskController) GetAllTasks(context *gin.Context) {

	tasks := ctrl.taskService.GetAllTasks()
	context.IndentedJSON(http.StatusOK, tasks)
}

func (ctrl *TaskController) AddTasks(context *gin.Context) {
	var payload models.Task
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	newTask, err := ctrl.taskService.AddTasks(payload)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Failed to create task"})
		return
	}
	context.IndentedJSON(http.StatusCreated, newTask)
}

func (ctrl *TaskController) DeleteTask(context *gin.Context) {
	idTask := context.Param("id")

	err := ctrl.taskService.DeleteTask(idTask)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found to delete"})
		return
	}
	context.Done()
}

func (ctrl *TaskController) GetTaskById(context *gin.Context) {
	id := context.Param("id")
	task, err := ctrl.taskService.FindTaskById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, task)
}

func (ctrl *TaskController) UpdateTask(context *gin.Context) {
	var input models.UpdateTaskInput
	idTask := context.Param("id")
	err := context.ShouldBind(&input)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	updateTask, err := ctrl.taskService.UpdateTask(idTask, input)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found to Update"})
		return
	}
	context.IndentedJSON(http.StatusOK, updateTask)
}
