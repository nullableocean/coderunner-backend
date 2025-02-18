package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"nullableocean-postupashki/src/models"
	"nullableocean-postupashki/src/service"

	"github.com/go-chi/chi/v5"
)

type CompileHandler struct {
	taskManager *service.CompileTaskManager
}

func NewTaskHandler(s *service.CompileTaskManager) *CompileHandler {
	return &CompileHandler{
		taskManager: s,
	}
}

type CreateResponse struct {
	TaskId string `json:"task_id"`
}

// @Summary Create task
// @Description "Create task for compile and execute. Need code and compiler name."
// @Tags Tasks
// @Accept json
// @Produce json
// @Param data body service.TaskCreateData true "data for execute"
// @Success 201 {object} CreateResponse "task uuid"
// @Failure default {string} string "error message"
// @Router /task [post]
func (h *CompileHandler) Create(w http.ResponseWriter, r *http.Request) {
	taskCreateData := service.TaskCreateData{}

	if err := json.NewDecoder(r.Body).Decode(&taskCreateData); err != nil {
		h.errorResponse(w, "invalid data", http.StatusBadRequest)
		return
	}

	if err := h.validateTaskData(&taskCreateData); err != nil {
		h.errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.taskManager.CreateTask(&taskCreateData)
	if err != nil {
		h.errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CreateResponse{TaskId: task.Uuid})
}

type StatusReponse struct {
	Status string `json:"status"`
}

// @Summary Get status
// @Description "Get task proccess status"
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} StatusReponse
// @Failure default {string} string "error message"
// @Router /status/{task_id} [get]
func (h *CompileHandler) Status(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "task_id")
	task, err := h.taskManager.GetTask(taskId)
	if err != nil {
		h.errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := h.taskManager.CheckTaskStatus(task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(StatusReponse{Status: status.String()})
}

type ResultResponse struct {
	Result string `json:"result"`
}

// @Summary Get result
// @Description "Get task execute result"
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} ResultResponse
// @Failure default {string} string "error message"
// @Router /result/{task_id} [get]
func (h *CompileHandler) Result(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "task_id")
	task, err := h.taskManager.GetTask(taskId)
	if err != nil {
		h.errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultRes := ResultResponse{Result: "wow result"}
	status := h.taskManager.CheckTaskStatus(task)
	if status == models.InProgress {
		resultRes.Result = status.String()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resultRes)
}

func (h *CompileHandler) errorResponse(w http.ResponseWriter, errMsg string, httpCode int) {
	http.Error(w, errMsg, httpCode)
}

func (h *CompileHandler) validateTaskData(data *service.TaskCreateData) error {
	if data.Code == "" {
		return errors.New("empty code")
	}

	if data.CompilerName == "" {
		return errors.New("empty compiler name")
	}

	return nil
}
