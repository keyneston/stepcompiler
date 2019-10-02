package step

import (
	"encoding/json"
	"time"
)

// Task is a state that does something. There are specific
// helper definitions for most common sub-types of a task, such as
// getting/putting/deleting from dynamodb or invoking a lambda and
// waiting for a callback.
//
// See https://docs.aws.amazon.com/step-functions/latest/dg/amazon-states-language-task-state.html
// for more details.
type Task struct {
	catch      []*CatchClause
	comment    string
	heartbeat  time.Duration
	next       State
	parameters map[string]interface{}
	resource   string
	timeout    time.Duration
	name       string
}

func (self *Task) Catch(input ...*CatchClause) *Task {
	self.catch = append(self.catch, input...)
	return self
}
func (self *Task) ChainableNext(input State) {
	self.Next(input)
}
func (self *Task) Comment(input string) *Task {
	self.comment = input
	return self
}

// Heartbeat is the number of seconds required between check-ins.
// If this time elapses without a check-in then the task is considered
// failed.
//
// Any time less than one second will induce a panic.
func (self *Task) Heartbeat(input time.Duration) *Task {
	self.heartbeat = input
	return self
}
func (self *Task) Next(input State) *Task {
	self.next = input
	return self
}
func (self *Task) Parameters(input map[string]interface{}) *Task {
	self.parameters = input
	return self
}
func (self *Task) Resource(input string) *Task {
	self.resource = input
	return self
}

// Timeout is the number of seconds for the task to complete.  If this
// time elapses without a check-in then the task is considered failed.
//
// Any time less than one second will induce a panic.
func (self *Task) Timeout(input time.Duration) *Task {
	self.timeout = input
	return self
}
func NewTask(name string) *Task {
	return &Task{name: name}
}
func (self Task) Name() string {
	return self.name
}
func (self Task) MarshalJSON() ([]byte, error) {
	out := &taskOutput{
		Catch:      self.catch,
		Comment:    self.comment,
		Heartbeat:  Timeout(self.heartbeat),
		Next:       "",
		Parameters: self.parameters,
		Resource:   self.resource,
		Timeout:    Timeout(self.timeout),
		Type:       self.StateType(),
	}
	if self.next != nil {
		out.Next = self.next.Name()
	} else {
		out.End = true
	}
	return json.Marshal(out)
}
func (self *Task) SetParameter(key string, value interface{}) *Task {
	if self.parameters == nil {
		self.parameters = map[string]interface{}{}
	}
	self.parameters[key] = value
	return self
}
func (self Task) GatherStates() []State {
	states := []State{self}
	if self.next != nil {
		states = append(states, self.next.GatherStates()...)
	}
	for _, clause := range self.catch {
		if clause.next != nil {
			states = append(states, clause.next.GatherStates()...)
		}
	}
	return states
}

type taskOutput struct {
	Catch      []*CatchClause         `json:"Catch,omitempty"`
	Comment    string                 `json:"Comment,omitempty"`
	End        bool                   `json:"End,omitempty"`
	Heartbeat  Timeout                `json:"HeartbeatSeconds,omitempty"`
	Next       string                 `json:"Next,omitempty"`
	Parameters map[string]interface{} `json:"Parameters,omitempty"`
	Resource   string                 `json:"Resource,omitempty"`
	Timeout    Timeout                `json:"TimeoutSeconds,omitempty"`
	Type       StateType              `json:"Type,omitempty"`
}

func (self Task) StateType() StateType {
	return StateType("Task")
}
