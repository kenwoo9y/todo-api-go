package entity

import "time"

type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "ToDo"
	TaskStatusDoing TaskStatus = "Doing"
	TaskStatusDone  TaskStatus = "Done"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     string     `json:"due_date"`
	Status      TaskStatus `json:"status"`
	OwnerID     int64      `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
