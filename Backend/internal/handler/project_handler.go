package handler

import (
	"encoding/json"
	"net/http"
	"taskflow/internal/middleware"
	service "taskflow/internal/services"
	"taskflow/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := utils.GetUserID(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	if req.Name == "" {
		utils.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	err := h.service.CreateProject(r.Context(), req.Name, req.Description, userID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]string{
		"message": "project created",
	})
}

func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {

	userID, ok := utils.GetUserID(r)
	if !ok {
		utils.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit, offset := utils.GetPagination(r)

	projects, err := h.service.GetProjects(r.Context(), userID, limit, offset)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to fetch projects")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"projects": projects,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "projectID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid project id")
		return
	}

	project, err := h.service.GetProjectByID(r.Context(), id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "not found")
		return
	}

	utils.JSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	var req CreateProjectRequest
	json.NewDecoder(r.Body).Decode(&req)

	claims := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	userID, _ := uuid.Parse(claims.UserID)

	err = h.service.UpdateProject(r.Context(), id, req.Name, req.Description, userID)
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

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(*utils.JWTClaims)
	userID, _ := uuid.Parse(claims.UserID)

	err = h.service.DeleteProject(r.Context(), id, userID)
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
