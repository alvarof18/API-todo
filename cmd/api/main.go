package main

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"time"
	"todo-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var tasks = []models.Task{}

func main() {

	//Iniciar servidor

	log.Println("Servidor inicializado")
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/task/:id", getTaskById)
	router.POST("/task", addTask)
	router.DELETE("/task/:id", deleteTask)
	router.PATCH("/task/:id", updateTask)
	router.Run("localhost:8080")

}

func getTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tasks)
}

func addTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	newTask.ID = uuid.New().String()
	newTask.CreatedAt = time.Now()
	newTask.Completed = false
	tasks = append(tasks, newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func getTaskById(context *gin.Context) {
	id := context.Param("id")

	for _, task := range tasks {
		if task.ID == id {
			context.IndentedJSON(http.StatusOK, task)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found"})
}

func deleteTask(context *gin.Context) {
	idTask := context.Param("id")
	task, err := getIndexTask(idTask)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found to delete"})
		return
	}

	tasks = append(tasks[:task], tasks[task+1:]...)
	context.Done()
}

func updateTask(context *gin.Context) {

	var input models.UpdateTaskInput

	idTask := context.Param("id")
	index, err := getIndexTask(idTask)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found to update"})
		return
	}

	err = context.ShouldBind(&input)

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// //Actualizo Informacion sin reflect
	// if input.Title != nil {
	// 	tasks[index].Title = *input.Title
	// }
	// if input.Completed != nil {
	// 	tasks[index].Completed = *input.Completed
	// }
	// context.IndentedJSON(http.StatusOK, tasks[index])
	//oriTask = &tasks[index]
	updateStruct(&tasks[index], &input)
	context.IndentedJSON(http.StatusOK, tasks[index])

}

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
