package step

import (
	"encoding/json"
	"time"
)

// DynamoGet is a wrapper for the DynamoDB getItem integration.
//
// AWS Documentation https://docs.aws.amazon.com/step-functions/latest/dg/connect-ddb.html
type DynamoGet struct {
	catch      []*CatchClause
	comment    string
	heartbeat  time.Duration
	next       State
	parameters map[string]interface{}
	timeout    time.Duration
	name       string
}

func (self *DynamoGet) Catch(input ...*CatchClause) *DynamoGet {
	self.catch = append(self.catch, input...)
	return self
}
func (self *DynamoGet) ChainableNext(input State) {
	self.Next(input)
}
func (self *DynamoGet) Comment(input string) *DynamoGet {
	self.comment = input
	return self
}

// Heartbeat is the number of seconds required between check-ins.
// If this time elapses without a check-in then the task is considered
// failed.
//
// Any time less than one second will induce a panic.
func (self *DynamoGet) Heartbeat(input time.Duration) *DynamoGet {
	self.heartbeat = input
	return self
}
func (self *DynamoGet) Key(input string) *DynamoGet {
	self.SetParameter("Key", input)
	return self
}
func (self *DynamoGet) Next(input State) *DynamoGet {
	self.next = input
	return self
}
func (self *DynamoGet) Parameters(input map[string]interface{}) *DynamoGet {
	self.parameters = input
	return self
}

// TableName sets the name of the table to make the dynamodb request to.
func (self *DynamoGet) TableName(input string) *DynamoGet {
	self.SetParameter("TableName", input)
	return self
}

// Timeout is the number of seconds for the task to complete.  If this
// time elapses without a check-in then the task is considered failed.
//
// Any time less than one second will induce a panic.
func (self *DynamoGet) Timeout(input time.Duration) *DynamoGet {
	self.timeout = input
	return self
}
func NewDynamoGet(name string) *DynamoGet {
	return &DynamoGet{name: name}
}
func (self DynamoGet) Name() string {
	return self.name
}
func (self DynamoGet) MarshalJSON() ([]byte, error) {
	out := &dynamogetOutput{
		Catch:      self.catch,
		Comment:    self.comment,
		Heartbeat:  Timeout(self.heartbeat),
		Next:       "",
		Parameters: self.parameters,
		Resource:   "arn:aws:states:::dynamodb:getItem",
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
func (self *DynamoGet) SetParameter(key string, value interface{}) *DynamoGet {
	if self.parameters == nil {
		self.parameters = map[string]interface{}{}
	}
	self.parameters[key] = value
	return self
}
func (self DynamoGet) GatherStates() []State {
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

type dynamogetOutput struct {
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

func (self DynamoGet) StateType() StateType {
	return StateType("Task")
}
