package task

import (
	"encoding/json"
	"errors"
	"io"
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

	err := decodeJSONBody(r, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.CreateTask(req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {

	filter := TaskFilter{
		Status:   strings.TrimSpace(r.URL.Query().Get("status")),
		Category: strings.TrimSpace(r.URL.Query().Get("category")),
		Search:   strings.TrimSpace(r.URL.Query().Get("search")),
	}

	tasks, err := h.service.GetTasks(filter)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) RemoveTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.DeleteTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "task deleted successfully",
	})
}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req UpdateTaskRequest
	err = decodeJSONBody(r, &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.UpdateTask(id, req)
	if err != nil {
		writeServiceError(w, err)
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

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, ErrorResponse{
		Error: message,
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTaskNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, ErrTitleRequired),
		errors.Is(err, ErrTitleTooShort),
		errors.Is(err, ErrTitleTooLong),
		errors.Is(err, ErrDescriptionTooLong),
		errors.Is(err, ErrCategoryTooLong),
		errors.Is(err, ErrInvalidStatus),
		errors.Is(err, ErrInvalidSearchFilter):
		writeError(w, http.StatusBadRequest, err.Error())

	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func decodeJSONBody(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return ErrInvalidRequestBody
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return ErrInvalidRequestBody
	}

	return nil
}

func getIDFromRequest(r *http.Request) (int, error) {
	idText := strings.TrimSpace(r.PathValue("id"))

	id, err := strconv.Atoi(idText)
	if err != nil {
		return 0, ErrInvalidTaskID
	}

	if id <= 0 {
		return 0, ErrInvalidTaskID
	}

	return id, nil
}
