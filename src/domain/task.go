package domain

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

type Task struct {
	Uuid         string     `json:"uuid"`
	Code         string     `json:"code"`
	CompilerName string     `json:"compiler"`
	Status       TaskStatus `json:"-"`
}
