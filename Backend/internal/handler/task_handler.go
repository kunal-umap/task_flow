package handler

import (
	"encoding/json"
	"net/http"
	"taskflow/internal/middleware"
	"taskflow/internal/models"
	service "taskflow/internal/services"
	"taskflow/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

type CreateTaskRequest struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Priority    models.TaskPriority `json:"priority"`
	AssigneeID  *string             `json:"assignee_id"`
	DueDate     *string             `json:"due_date"` // ISO string
}

type UpdateTaskRequest struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Status      models.TaskStatus   `json:"status"`
	Priority    models.TaskPriority `json:"priority"`
	AssigneeID  *string             `json:"assignee_id"`
	DueDate     *string             `json:"due_date"`
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {

	taskIDStr := chi.URLParam(r, "taskID")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	userID, ok := utils.GetUserID(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), taskID, userID)
	if err != nil {
		if err.Error() == "task not found" {
			utils.Error(w, http.StatusNotFound, "not found")
			return
		}
		if err.Error() == "forbidden" {
			utils.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "error")
		return
	}

	utils.JSON(w, http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var assigneeUUID *uuid.UUID
	if req.AssigneeID != nil {
		id, err := uuid.Parse(*req.AssigneeID)
		if err == nil {
			assigneeUUID = &id
		}
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err == nil {
			dueDate = &t
		}
	}

	err = h.service.CreateTask(
		r.Context(),
		req.Title,
		req.Description,
		projectID,
		assigneeUUID,
		req.Priority,
		dueDate,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {

	projectIDStr := chi.URLParam(r, "projectID")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid project id")
		return
	}

	limit, offset := utils.GetPagination(r)

	var status *models.TaskStatus
	if s := r.URL.Query().Get("status"); s != "" {
		st := models.TaskStatus(s)
		status = &st
	}

	var assigneeID *uuid.UUID
	if a := r.URL.Query().Get("assignee"); a != "" {
		id, _ := uuid.Parse(a)
		assigneeID = &id
	}

	tasks, err := h.service.GetTasks(r.Context(), projectID, status, assigneeID, limit, offset)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"tasks":  tasks,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "taskID")
	taskID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	json.NewDecoder(r.Body).Decode(&req)

	claims := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	userID, _ := uuid.Parse(claims.UserID)

	var assigneeUUID *uuid.UUID
	if req.AssigneeID != nil {
		id, _ := uuid.Parse(*req.AssigneeID)
		assigneeUUID = &id
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.DueDate)
		dueDate = &t
	}

	err = h.service.UpdateTask(
		r.Context(),
		taskID,
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		assigneeUUID,
		dueDate,
		userID,
	)

	if err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("updated"))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "taskID")

	taskID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	userID, _ := uuid.Parse(claims.UserID)

	err = h.service.DeleteTask(r.Context(), taskID, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
