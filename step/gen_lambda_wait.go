package step

import (
	"encoding/json"
	"time"
)

type LambdaWait struct {
	catch      []*CatchClause
	comment    string
	heartbeat  time.Duration
	next       State
	parameters map[string]interface{}
	resource   string
	timeout    time.Duration
	name       string
}

func (self *LambdaWait) Catch(input ...*CatchClause) *LambdaWait {
	self.catch = append(self.catch, input...)
	return self
}
func (self *LambdaWait) ChainableNext(input State) {
	self.Next(input)
}
func (self *LambdaWait) Comment(input string) *LambdaWait {
	self.comment = input
	return self
}

// FunctionName sets the ARN/Name/or other indicator of the Lambda
// Function to invoke. See AWS documentation for more details:
// https://docs.aws.amazon.com/lambda/latest/dg/API_Invoke.html#API_Invoke_RequestParameters
func (self *LambdaWait) FunctionName(input string) *LambdaWait {
	self.SetParameter("FunctionName", input)
	return self
}

// Heartbeat is the number of seconds required between check-ins.
// If this time elapses without a check-in then the task is considered
// failed.
//
// Any time less than one second will induce a panic.
func (self *LambdaWait) Heartbeat(input time.Duration) *LambdaWait {
	self.heartbeat = input
	return self
}
func (self *LambdaWait) Next(input State) *LambdaWait {
	self.next = input
	return self
}
func (self *LambdaWait) Parameters(input map[string]interface{}) *LambdaWait {
	self.parameters = input
	return self
}
func (self *LambdaWait) Resource(input string) *LambdaWait {
	self.resource = input
	return self
}

// Timeout is the number of seconds for the task to complete.  If this
// time elapses without a check-in then the task is considered failed.
//
// Any time less than one second will induce a panic.
func (self *LambdaWait) Timeout(input time.Duration) *LambdaWait {
	self.timeout = input
	return self
}
func NewLambdaWait(name string) *LambdaWait {
	return &LambdaWait{name: name}
}
func (self LambdaWait) Name() string {
	return self.name
}
func (self LambdaWait) MarshalJSON() ([]byte, error) {
	out := &lambdawaitOutput{
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
func (self *LambdaWait) SetParameter(key string, value interface{}) *LambdaWait {
	if self.parameters == nil {
		self.parameters = map[string]interface{}{}
	}
	self.parameters[key] = value
	return self
}
func (self LambdaWait) GatherStates() []State {
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

type lambdawaitOutput struct {
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

func (self LambdaWait) StateType() StateType {
	return StateType("Task")
}
