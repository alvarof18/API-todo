package models

import (
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" binding:"required"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createAt"`
}

// Se define como puntero ya que en caso de no enviarse me llega como null
type UpdateTaskInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}
