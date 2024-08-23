package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"todo-api/internal/models"
)

type TaskRepository interface {
	GetAllTasks() []models.Task
	FindTaskById(id string) (*models.Task, error)
	AddTasks(payload models.Task) error
	DeleteTask(id string) error
	UpdateTask(newTask models.Task) error
}

type taskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskReposiory(db *sql.DB) taskRepositoryImpl {
	return taskRepositoryImpl{db: db}
}

// AddTasks implements TaskRepository.
func (t *taskRepositoryImpl) AddTasks(newTask models.Task) error {
	_, err := t.db.Exec("INSERT INTO tasks VALUES(?,?,?,?)", newTask.ID, newTask.Title, newTask.Completed, newTask.CreatedAt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
}

// DeleteTask implements TaskRepository.
func (t *taskRepositoryImpl) DeleteTask(id string) error {
	result, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("task not found to delete")
	}
	fmt.Printf("Rows affected %v", rows)
	return nil
}

// FindTaskById implements TaskRepository.
func (t *taskRepositoryImpl) FindTaskById(id string) (*models.Task, error) {
	var task models.Task
	result := t.db.QueryRow("SELECT * FROM tasks where id =?", id)

	if err := result.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %q: no found", id)
		}
		return nil, fmt.Errorf(`error %v`, err)
	}
	return &task, nil
}

// GetAllTasks implements TaskRepository.
func (t *taskRepositoryImpl) GetAllTasks() []models.Task {
	var tasksResult = []models.Task{}

	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return tasksResult
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt); err != nil {
			fmt.Printf("Error Getting task %v", err)
			return tasksResult
		}
		tasksResult = append(tasksResult, task)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error Getting tasks %v", err)
		return tasksResult
	}
	return tasksResult
}

// UpdateTask implements TaskRepository.
func (t *taskRepositoryImpl) UpdateTask(newTask models.Task) error {
	_, err := t.db.Exec("UPDATE tasks SET title =?, completed =? WHERE id=?", newTask.Title, newTask.Completed, newTask.ID)
	if err != nil {
		fmt.Printf("Error Update task: %v\n", err)
		return err
	}
	return nil
}
