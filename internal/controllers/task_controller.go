package controllers

import (
	"net/http"
	"todo-api/internal/models"
	"todo-api/internal/services"

	"github.com/gin-gonic/gin"
)

// Maneja las peticiones HTTP
type TaskController struct{}

var servicesTask = services.TaskService{}

func (c *TaskController) GetAllTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, servicesTask.GetAllTasks())
}

func (c *TaskController) AddTasks(context *gin.Context) {
	var payload models.Task
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	newTask, err := servicesTask.AddTasks(payload)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Error": "Failed to create task"})
		return
	}
	context.IndentedJSON(http.StatusOK, newTask)
}

func (c *TaskController) DeleteTask(context *gin.Context) {
	idTask := context.Param("id")

	err := servicesTask.DeleteTask(idTask)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found to delete"})
		return
	}
	context.Done()
}

func (c *TaskController) GetTaskById(context *gin.Context) {
	id := context.Param("id")
	task, err := servicesTask.FindTaskById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, task)
}

func (c *TaskController) UpdateTask(context *gin.Context) {
	var input models.UpdateTaskInput
	idTask := context.Param("id")
	err := context.ShouldBind(&input)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	updateTask, err := servicesTask.UpdateTask(idTask, input)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found to Update"})
		return
	}
	context.IndentedJSON(http.StatusOK, updateTask)
}
