package stepcompiler

const TaskType = "Task"

type Task struct {
	name       string                 `json:"-"`
	Type       string                 `json:"Type"`
	Parameters map[string]interface{} `json:"Parameters,omitempty"`
	ResultPath string                 `json:"ResultPath,omitempty"`
	Resource   string                 `json:"Resource,omitempty"`
	next       State                  `json:"Next,omitempty"`
	Catch      []CatchClause          `json:"Catch,omitempty"`
	Comment    string                 `json:"Comment,omitempty"`
}

func (Task) StateType() string {
	return TaskType
}

func (t Task) Name() string {
	return t.name
}

func NewTask(name string) *Task {
	return &Task{
		name: name,
		Type: TaskType,
	}
}

func (t *Task) SetComment(comment string) *Task {
	t.Comment = comment
	return t
}

func (t *Task) Next(state State) *Task {
	t.next = state

	return t
}

func (t *Task) GatherStates() []State {
	res := []State{}

	if t.next != nil {
		res = append(res, t.next)
	}

	return res
}
