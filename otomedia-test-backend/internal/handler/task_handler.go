package handler

import (
	"net/http"
	"strconv"

	"otomedia/task-managment/internal/dto"
	"otomedia/task-managment/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// List menangani GET /api/tasks?status=&keyword=&assignee=&page=&limit=&sort=
func (h *TaskHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	sort := c.DefaultQuery("sort", "newest")
	if sort != "newest" && sort != "oldest" {
		sort = "newest"
	}

	status := c.Query("status")
	// "all" dipakai frontend untuk merepresentasikan "tanpa filter status"
	if status == "all" {
		status = ""
	}

	filter := dto.ListFilter{
		Status:   status,
		Keyword:  c.Query("keyword"),
		Assignee: c.Query("assignee"),
		Page:     page,
		Limit:    limit,
		Sort:     sort,
	}

	resp, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Create menangani POST /api/tasks
func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, err.Error())
		return
	}

	task, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": task})
}

// Update menangani PUT /api/tasks/{id}
func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondValidationError(c, "invalid task id")
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondValidationError(c, err.Error())
		return
	}

	task, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		RespondError(c, err) // duplicate title -> 409, not found -> 404 (Task 4)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

// Delete menangani DELETE /api/tasks/{id} (soft delete)
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		RespondValidationError(c, "invalid task id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		RespondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
