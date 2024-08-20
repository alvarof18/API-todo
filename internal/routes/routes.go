package routes

import (
	"todo-api/internal/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	taskController := controllers.TaskController{}

	router.GET("/tasks", taskController.GetAllTasks)
	router.GET("/task/:id", taskController.GetTaskById)
	router.POST("/task", taskController.AddTasks)
	router.DELETE("/task/:id", taskController.DeleteTask)
	router.PATCH("/task/:id", taskController.UpdateTask)
}
