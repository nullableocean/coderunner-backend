package rest

import (
	"net/http"
	"nullableocean-postupashki/src/api/rest/middleware"
	"nullableocean-postupashki/src/api/rest/types"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/usecases"

	"github.com/go-chi/chi/v5"
)

type CommitHandler struct {
	resultService usecases.Result

	accessHeader string
	accessToken  string
}

func NewCommitHandler(accessHeader, accessToken string, rs usecases.Result) *CommitHandler {
	return &CommitHandler{
		resultService: rs,
		accessHeader:  accessHeader,
		accessToken:   accessToken,
	}
}

func (h *CommitHandler) RegisterRoutes(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.CommiterAuth(h.accessHeader, h.accessToken)) //secure

		r.Post("/commit", h.Post)
	})
}

// @Summary Commit result
// @Description "Save result"
// @Tags Results
// @Secure APIKeyHeader
// @Accept json
// @Produce json
// @Param data body types.PostCommitRequestBody true "data for save"
// @Success 200
// @Failure 401
// @Failure default {object} types.ErrorResponse "error response"
// @Router /commit [post]
func (h *CommitHandler) Post(w http.ResponseWriter, r *http.Request) {
	resultBody, err := types.ExtractPostResultBody(r)
	if err != nil {
		types.ProcessError(w, err, http.StatusBadRequest)
		return
	}

	res := &domain.Result{
		TaskUuid:  resultBody.TaskUuid,
		Output:    resultBody.Output,
		Error:     resultBody.Error,
		IsSuccess: resultBody.IsSuccess,
	}

	err = h.resultService.SaveResult(res)
	if err != nil {
		types.ProcessError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
