package stepcompiler

import "encoding/json"

const TaskType = "Task"

type Task struct {
	name       string
	parameters map[string]interface{} `json:"Parameters,omitempty"`
	resultPath string                 `json:"ResultPath,omitempty"`
	resource   string                 `json:"Resource,omitempty"`
	next       State                  `json:"Next,omitempty"`
	catch      []CatchClause          `json:"Catch,omitempty"`
	comment    string                 `json:"Comment,omitempty"`
}

type taskOutput struct {
	Type       string                 `json:"Type"`
	Parameters map[string]interface{} `json:"Parameters,omitempty"`
	ResultPath string                 `json:"ResultPath,omitempty"`
	Resource   string                 `json:"Resource,omitempty"`
	Next       string                 `json:"Next,omitempty"`
	Catch      []CatchClause          `json:"Catch,omitempty"`
	Comment    string                 `json:"Comment,omitempty"`
	End        bool                   `json:"End,omitempty"`
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
	}
}

func (t *Task) Comment(comment string) *Task {
	t.comment = comment
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

func (t Task) MarshalJSON() ([]byte, error) {
	out := taskOutput{
		Type:    t.StateType(),
		Comment: t.comment,
	}

	if t.next != nil {
		out.Next = t.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}
