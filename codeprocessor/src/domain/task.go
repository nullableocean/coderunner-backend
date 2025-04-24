package domain

type Task struct {
	Uuid         string `json:"uuid"`
	Code         string `json:"code"`
	CompilerName string `json:"compiler"`
}
