package models

type TaskStatus int

const (
	Ready TaskStatus = iota
	InProgress
)

func (ts TaskStatus) String() string {
	switch ts {
	case Ready:
		return "ready"
	case InProgress:
		return "in_progress"
	}

	return ""
}

type CompileTask struct {
	Uuid         string
	Code         string
	CompilerName string
	Status       TaskStatus
}
