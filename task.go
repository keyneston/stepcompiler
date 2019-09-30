package stepcompiler

import (
	"encoding/json"
	"time"
)

// Task is a builder for building Tasks.
//
// See https://docs.aws.amazon.com/step-functions/latest/dg/amazon-states-language-task-state.html for full details.
type Task struct {
	name       string
	parameters map[string]interface{}
	resultPath string
	resource   string
	next       State
	catch      []*CatchClause
	comment    string
	timeout    time.Duration
}

type taskOutput struct {
	Type             StateType              `json:"Type"`
	HeartbeatSeconds Timeout                `json:"HeartbeatSeconds,omitempty"`
	Parameters       map[string]interface{} `json:"Parameters,omitempty"`
	ResultPath       string                 `json:"ResultPath,omitempty"`
	Resource         string                 `json:"Resource,omitempty"`
	Next             string                 `json:"Next,omitempty"`
	Catch            []*CatchClause         `json:"Catch,omitempty"`
	Comment          string                 `json:"Comment,omitempty"`
	End              bool                   `json:"End,omitempty"`
	TimeoutSeconds   Timeout                `json:"TimeoutSeconds,omitempty"`
}

func (Task) StateType() StateType {
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

func (t *Task) Parameters(params map[string]interface{}) *Task {
	t.parameters = params

	return t
}

func (t *Task) ResultPath(resultPath string) *Task {
	t.resultPath = resultPath
	return t
}

func (t *Task) Resource(resource string) *Task {
	t.resource = resource
	return t
}

func (t *Task) Comment(comment string) *Task {
	t.comment = comment
	return t
}

func (t *Task) Catch(clause *CatchClause) *Task {
	t.catch = append(t.catch, clause)

	return t
}

func (t *Task) CatchChainable(clause *CatchClause) {
	t.Catch(clause)
}

func (t *Task) Next(state State) *Task {
	t.next = state

	return t
}

func (t *Task) NextChainable(state State) {
	t.Next(state)
}

func (t *Task) GatherStates() []State {
	res := []State{t}

	if t.next != nil {
		res = append(res, t.next)
		res = append(res, t.next.GatherStates()...)
	}

	for _, clause := range t.catch {
		res = append(res, clause.GatherStates()...)
	}

	return res
}

func (t *Task) Timeout(timeout time.Duration) *Task {
	t.timeout = timeout
	return t
}

func (t Task) MarshalJSON() ([]byte, error) {
	out := taskOutput{
		Type:           t.StateType(),
		Comment:        t.comment,
		Resource:       t.resource,
		Parameters:     t.parameters,
		ResultPath:     t.resultPath,
		Catch:          t.catch,
		TimeoutSeconds: Timeout(t.timeout),
	}

	if t.next != nil {
		out.Next = t.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}
