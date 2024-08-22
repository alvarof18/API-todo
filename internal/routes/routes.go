package routes

import (
	"todo-api/internal/controllers"
	"todo-api/internal/services"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	taskService := &services.TaskServiceImpl{}
	taskController := controllers.NewTaskController(taskService)

	router.GET("/tasks", taskController.GetAllTasks)
	router.GET("/task/:id", taskController.GetTaskById)
	router.POST("/task", taskController.AddTasks)
	router.DELETE("/task/:id", taskController.DeleteTask)
	router.PATCH("/task/:id", taskController.UpdateTask)
}
