package rest

import (
	"encoding/json"
	"net/http"
	"nullableocean-postupashki/src/api/rest/types"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/usecases/service"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService *service.CompileTaskService
}

func NewTaskHandler(s *service.CompileTaskService) *TaskHandler {
	return &TaskHandler{
		taskService: s,
	}
}

// @Summary Create task and process
// @Description "Create task for compile and execute. Need code and compiler name."
// @Tags Tasks
// @Accept json
// @Produce json
// @Param data body types.PostTaskRequestBody true "data for execute"
// @Success 201 {object} types.PostTaskResponse "task uuid"
// @Failure default {object} types.ErrorResponse "error response"
// @Router /task [post]
func (h *TaskHandler) Post(w http.ResponseWriter, r *http.Request) {
	taskPostBody, err := types.ExtractPostTaskBody(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.taskService.CreateAndProcess(&service.TaskCreateData{
		Code:         taskPostBody.Code,
		CompilerName: taskPostBody.CompilerName,
	})

	if err != nil {
		types.ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&types.PostTaskResponse{TaskId: task.Uuid})
}

// @Summary Get GetStatus
// @Description "Get task process GetStatus"
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} types.StatusResponse
// @Failure default {object} types.ErrorResponse "error response"
// @Router /status/{task_id} [get]
func (h *TaskHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	task, err := h.getTaskFromRequest(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	status := h.taskService.CheckStatus(task)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&types.StatusResponse{Status: status.String()})
}

// @Summary Get result
// @Description "Get task execute result"
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} types.ResultResponse
// @Failure default {object} types.ErrorResponse "error response"
// @Router /result/{task_id} [get]
func (h *TaskHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	task, err := h.getTaskFromRequest(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	resultRes := types.ResultResponse{Result: "wow result"}

	status := h.taskService.CheckStatus(task)
	if status == domain.InProgress {
		resultRes.Result = status.String()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resultRes)
}

func (h *TaskHandler) getTaskFromRequest(r *http.Request) (*domain.Task, error) {
	taskId := chi.URLParam(r, "task_id")
	return h.taskService.GetTask(taskId)
}

func (h *TaskHandler) RegisterRoutes(router *chi.Mux) {
	router.Post("/task", h.Post)
	router.Get("/status/{task_id}", h.GetStatus)
	router.Get("/result/{task_id}", h.GetResult)
}
