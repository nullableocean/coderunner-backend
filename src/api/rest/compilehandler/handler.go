package compilehandler

import (
	"encoding/json"
	"net/http"
	"nullableocean-postupashki/src/models"
	"nullableocean-postupashki/src/service/compile"

	"github.com/go-chi/chi/v5"
)

type CompileHandler struct {
	compileService *compile.CompileTaskService
}

func NewCompileHandler(s *compile.CompileTaskService) *CompileHandler {
	return &CompileHandler{
		compileService: s,
	}
}

// @Summary Create task
// @Description "Create task for compile and execute. Need code and compiler name."
// @Tags Tasks
// @Accept json
// @Produce json
// @Param data body CreateTaskRequestBody true "data for execute"
// @Success 201 {object} CreateTaskResponse "task uuid"
// @Failure default {object} rest.ErrorResponse "error response"
// @Router /task [post]
func (h *CompileHandler) create(w http.ResponseWriter, r *http.Request) {
	taskCreateBody, err := extractCreateTaskBody(r)
	if err != nil {
		proccessError(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.compileService.CreateTask(&compile.TaskCreateData{
		Code:         taskCreateBody.Code,
		CompilerName: taskCreateBody.CompilerName,
	})

	if err != nil {
		proccessError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&CreateTaskResponse{TaskId: task.Uuid})
}

// @Summary Get status
// @Description "Get task proccess status"
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} StatusResponse
// @Failure default {object} rest.ErrorResponse "error response"
// @Router /status/{task_id} [get]
func (h *CompileHandler) status(w http.ResponseWriter, r *http.Request) {
	task, err := h.getTaskFromRequest(r)
	if err != nil {
		proccessError(w, err, http.StatusInternalServerError)
		return
	}

	status := h.compileService.CheckTaskStatus(task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&StatusResponse{Status: status.String()})
}

// @Summary Get result
// @Description "Get task execute result"
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} ResultResponse
// @Failure default {object} rest.ErrorResponse "error response"
// @Router /result/{task_id} [get]
func (h *CompileHandler) result(w http.ResponseWriter, r *http.Request) {
	task, err := h.getTaskFromRequest(r)
	if err != nil {
		proccessError(w, err, http.StatusInternalServerError)
		return
	}

	resultRes := ResultResponse{Result: "wow result"}

	status := h.compileService.CheckTaskStatus(task)
	if status == models.InProgress {
		resultRes.Result = status.String()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resultRes)
}

func (h *CompileHandler) getTaskFromRequest(r *http.Request) (*models.CompileTask, error) {
	taskId := chi.URLParam(r, "task_id")
	return h.compileService.GetTask(taskId)
}

func (h *CompileHandler) RegisterRoutes(router *chi.Mux) {
	router.Post("/task", h.create)
	router.Get("/status/{task_id}", h.status)
	router.Get("/result/{task_id}", h.result)
}
