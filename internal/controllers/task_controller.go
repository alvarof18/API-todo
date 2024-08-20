package controllers

import (
	"errors"
	"net/http"
	"reflect"
	"time"
	"todo-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskController struct{}

// var tasks = []models.Task{
// 	{Title: "Comer", Completed: false, CreatedAt: time.Time{}},
// 	{Title: "Estudiar", Completed: false, CreatedAt: time.Time{}},
// 	{Title: "Programar", Completed: false, CreatedAt: time.Time{}},
// }

var tasks = []models.Task{}

func (c *TaskController) GetAllTasks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, tasks)
}

func (c *TaskController) AddTasks(context *gin.Context) {
	var newTask models.Task
	if err := context.ShouldBindJSON(&newTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	newTask.ID = uuid.New().String()
	newTask.CreatedAt = time.Now()
	newTask.Completed = false
	tasks = append(tasks, newTask)
	context.IndentedJSON(http.StatusCreated, newTask)
}

func (c *TaskController) DeleteTask(context *gin.Context) {
	idTask := context.Param("id")
	task, err := getIndexTask(idTask)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found to delete"})
		return
	}
	tasks = append(tasks[:task], tasks[task+1:]...)
	context.Done()
}

func (c *TaskController) GetTaskById(context *gin.Context) {
	id := context.Param("id")
	for _, task := range tasks {
		if task.ID == id {
			context.IndentedJSON(http.StatusOK, task)
			return
		}
	}
	context.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Task not found"})
}

func (c *TaskController) UpdateTask(context *gin.Context) {

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
	updateStruct(&tasks[index], &input)
	context.IndentedJSON(http.StatusOK, tasks[index])

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
