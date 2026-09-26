package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"otomedia/task-managment/internal/apperror"
	"otomedia/task-managment/internal/cache"
	"otomedia/task-managment/internal/dto"
	"otomedia/task-managment/internal/models"
	"otomedia/task-managment/internal/repository"
)

// listCacheTTL: GET /api/tasks di-cache 60 detik (Task 2).
const listCacheTTL = 60 * time.Second

type TaskService struct {
	repo  repository.TaskRepository
	cache cache.Cache
}

func NewTaskService(repo repository.TaskRepository, c cache.Cache) *TaskService {
	return &TaskService{repo: repo, cache: c}
}

// buildListCacheKey menyusun cache key yang mengikutsertakan seluruh query
// parameter (status, keyword, assignee, page, limit, sort) sesuai requirement
// Task 2: "Cache key must include query parameters".
func buildListCacheKey(f dto.ListFilter) string {
	return fmt.Sprintf(
		"tasks:list:status=%s:keyword=%s:assignee=%s:page=%d:limit=%d:sort=%s",
		f.Status, f.Keyword, f.Assignee, f.Page, f.Limit, f.Sort,
	)
}

func (s *TaskService) List(ctx context.Context, filter dto.ListFilter) (*dto.TaskListResponse, error) {
	key := buildListCacheKey(filter)

	if cached, hit, err := s.cache.Get(ctx, key); err == nil && hit {
		var resp dto.TaskListResponse
		if jsonErr := json.Unmarshal([]byte(cached), &resp); jsonErr == nil {
			return &resp, nil
		}
		// cache korup/format berubah -> abaikan dan query ulang ke DB
	}

	tasks, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, apperror.Internal("failed to fetch tasks")
	}

	resp := &dto.TaskListResponse{Data: tasks, Page: filter.Page, Limit: filter.Limit, Total: total}

	if payload, err := json.Marshal(resp); err == nil {
		// Kegagalan menulis cache tidak boleh menggagalkan request GET,
		// jadi errornya sengaja diabaikan (best-effort caching).
		_ = s.cache.Set(ctx, key, string(payload), listCacheTTL)
		_ = s.cache.TrackListKey(ctx, key)
	}

	return resp, nil
}

func (s *TaskService) Create(ctx context.Context, req dto.CreateTaskRequest) (*models.Task, error) {
	if req.Status == "" {
		req.Status = models.StatusTodo
	}
	if !req.Status.Valid() {
		return nil, apperror.BadRequest("invalid status value")
	}

	existing, err := s.repo.FindActiveByTitle(ctx, req.Title, 0)
	if err != nil {
		return nil, apperror.Internal("failed to validate title")
	}
	if existing != nil {
		return nil, apperror.Conflict("a task with this title already exists")
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Assignee:    req.Assignee,
	}
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, apperror.Internal("failed to create task")
	}

	_ = s.cache.InvalidateListCache(ctx)
	return task, nil
}

func (s *TaskService) Update(ctx context.Context, id int64, req dto.UpdateTaskRequest) (*models.Task, error) {
	if !req.Status.Valid() {
		return nil, apperror.BadRequest("invalid status value")
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("failed to fetch task")
	}
	if current == nil {
		return nil, apperror.NotFound("task not found")
	}

	// Task 4 bug fix: duplicate title saat update harus 409, bukan 500.
	existing, err := s.repo.FindActiveByTitle(ctx, req.Title, id)
	if err != nil {
		return nil, apperror.Internal("failed to validate title")
	}
	if existing != nil {
		return nil, apperror.Conflict("a task with this title already exists")
	}

	current.Title = req.Title
	current.Description = req.Description
	current.Status = req.Status
	current.Assignee = req.Assignee

	if err := s.repo.Update(ctx, current); err != nil {
		return nil, apperror.Internal("failed to update task")
	}

	_ = s.cache.InvalidateListCache(ctx)
	return current, nil
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperror.Internal("failed to fetch task")
	}
	if current == nil {
		return apperror.NotFound("task not found")
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return apperror.Internal("failed to delete task")
	}

	_ = s.cache.InvalidateListCache(ctx)
	return nil
}
