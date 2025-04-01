package domain

type Result struct {
	TaskUuid  string `json:"task_uuid"`
	Output    string `json:"output"`
	Error     string `json:"error"`
	IsSuccess bool   `json:"is_success"`
}
