package rest

import (
	"encoding/json"
	"net/http"
	"nullableocean-postupashki/src/api/rest/middleware"
	"nullableocean-postupashki/src/api/rest/types"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/usecases"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService    usecases.Task
	resultService  usecases.Result
	sessionService usecases.Session
}

func NewTaskHandler(s usecases.Task, resultService usecases.Result, sessions usecases.Session) *TaskHandler {
	return &TaskHandler{
		taskService:    s,
		resultService:  resultService,
		sessionService: sessions,
	}
}

func (h *TaskHandler) RegisterRoutes(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth(h.sessionService)) //secure

		r.Post("/task", h.Post)
		r.Get("/status/{task_id}", h.GetStatus)
		r.Get("/result/{task_id}", h.GetResult)
	})
}

// @Summary Create task and process
// @Description "Create task for compile and execute. Need code and compiler name."
// @Tags Tasks
// @Secure APIKeyHeader
// @Accept json
// @Produce json
// @Param data body types.PostTaskRequestBody true "data for execute"
// @Success 201 {object} types.PostTaskResponse "task uuid"
// @Failure 401
// @Failure default {object} types.ErrorResponse "error response"
// @Router /task [post]
func (h *TaskHandler) Post(w http.ResponseWriter, r *http.Request) {
	taskPostBody, err := types.ExtractPostTaskBody(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.taskService.CreateAndProcess(&domain.Task{
		Code:         taskPostBody.Code,
		CompilerName: taskPostBody.CompilerName,
	})

	if err != nil {
		types.ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&types.PostTaskResponse{TaskId: task.Uuid})
}

// @Summary Get GetStatus
// @Description "Get task process GetStatus"
// @Secure APIKeyHeader
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} types.StatusResponse
// @Failure 401
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
// @Secure APIKeyHeader
// @Tags Tasks
// @Produce json
// @Param task_id path string true "task uuid"
// @Success 201 {object} types.ResultResponse
// @Failure 401
// @Failure default {object} types.ErrorResponse "error response"
// @Router /result/{task_id} [get]
func (h *TaskHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	task, err := h.getTaskFromRequest(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	resultRes := types.ResultResponse{}

	status := h.taskService.CheckStatus(task)
	if status == domain.Ready {
		result, err := h.resultService.GetTaskResult(task)
		if err != nil {
			types.ProcessError(w, err, http.StatusInternalServerError)
			return
		}
		resultRes.Result = result.Output
	} else {
		resultRes.Result = status.String()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resultRes)
}

func (h *TaskHandler) getTaskFromRequest(r *http.Request) (*domain.Task, error) {
	taskId := chi.URLParam(r, "task_id")
	return h.taskService.Get(taskId)
}
