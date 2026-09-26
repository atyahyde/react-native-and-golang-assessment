package dto

import "otomedia/task-managment/internal/models"

// ListFilter menampung semua query parameter yang didukung GET /api/tasks:
// status, keyword, assignee, page, limit, sort.
type ListFilter struct {
	Status   string
	Keyword  string
	Assignee string
	Page     int
	Limit    int
	Sort     string // "newest" | "oldest"
}

type CreateTaskRequest struct {
	Title       string            `json:"title" binding:"required,min=3,max=255"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
	Assignee    string            `json:"assignee"`
}

type UpdateTaskRequest struct {
	Title       string            `json:"title" binding:"required,min=3,max=255"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status" binding:"required"`
	Assignee    string            `json:"assignee"`
}

type TaskListResponse struct {
	Data  []models.Task `json:"data"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Total int64         `json:"total"`
}
