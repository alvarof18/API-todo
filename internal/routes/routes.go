package routes

import (
	"database/sql"
	"todo-api/internal/controllers"
	"todo-api/internal/services"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine, db *sql.DB) {

	taskService := services.NewTaskService(db)
	taskController := controllers.NewTaskController(taskService)

	router.GET("/tasks", taskController.GetAllTasks)
	router.GET("/task/:id", taskController.GetTaskById)
	router.POST("/task", taskController.AddTasks)
	router.DELETE("/task/:id", taskController.DeleteTask)
	router.PATCH("/task/:id", taskController.UpdateTask)
}
