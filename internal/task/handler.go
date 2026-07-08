package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	task, err := h.service.CreateTask(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.service.GetTasks()

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetIDFromPath(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) RemoveTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetIDFromPath(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}
	err = h.service.DeleteTask(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "task deleted successfully",
	})
}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.GetIDFromPath(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	var req UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	task, err := h.service.UpdateTask(id, req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) GetIDFromPath(path string) (int, error) {
	idText := strings.TrimPrefix(path, "/tasks/")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
