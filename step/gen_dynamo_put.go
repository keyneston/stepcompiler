package step

import (
	"encoding/json"
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

// DynamoPut is a wrapper for the DynamoDB putItem integration.
//
// AWS Documentation https://docs.aws.amazon.com/step-functions/latest/dg/connect-ddb.html
type DynamoPut struct {
	catch      []*CatchClause
	comment    string
	heartbeat  time.Duration
	next       State
	parameters map[string]interface{}
	resource   string
	resultpath string
	timeout    time.Duration
	name       string
}

func (self *DynamoPut) Catch(input ...*CatchClause) *DynamoPut {
	self.catch = append(self.catch, input...)
	return self
}
func (self *DynamoPut) ChainableNext(input State) {
	self.Next(input)
}
func (self *DynamoPut) Comment(input string) *DynamoPut {
	self.comment = input
	return self
}
func (self *DynamoPut) ConditionExpression(input string) *DynamoPut {
	self.SetParameter("ConditionExpression", input)
	return self
}

// Heartbeat is the number of seconds required between check-ins.
// If this time elapses without a check-in then the task is considered
// failed.
//
// Any time less than one second will induce a panic.
func (self *DynamoPut) Heartbeat(input time.Duration) *DynamoPut {
	self.heartbeat = input
	return self
}
func (self *DynamoPut) Item(input map[string]*dynamodb.AttributeValue) *DynamoPut {
	self.SetParameter("Key", input)
	return self
}
func (self *DynamoPut) Next(input State) *DynamoPut {
	self.next = input
	return self
}
func (self *DynamoPut) Parameters(input map[string]interface{}) *DynamoPut {
	self.parameters = input
	return self
}
func (self *DynamoPut) Resource(input string) *DynamoPut {
	self.resource = input
	return self
}
func (self *DynamoPut) ResultPath(input string) *DynamoPut {
	self.resultpath = input
	return self
}

// TableName sets the name of the table to make the dynamodb request to.
func (self *DynamoPut) TableName(input string) *DynamoPut {
	self.SetParameter("TableName", input)
	return self
}

// Timeout is the number of seconds for the task to complete.  If this
// time elapses without a check-in then the task is considered failed.
//
// Any time less than one second will induce a panic.
func (self *DynamoPut) Timeout(input time.Duration) *DynamoPut {
	self.timeout = input
	return self
}
func NewDynamoPut(name string) *DynamoPut {
	return &DynamoPut{name: name}
}
func (self DynamoPut) Name() string {
	return self.name
}
func (self DynamoPut) MarshalJSON() ([]byte, error) {
	out := &dynamoputOutput{
		Catch:      self.catch,
		Comment:    self.comment,
		Heartbeat:  Timeout(self.heartbeat),
		Next:       "",
		Parameters: self.parameters,
		Resource:   self.resource,
		ResultPath: self.resultpath,
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
func (self *DynamoPut) SetParameter(key string, value interface{}) *DynamoPut {
	if self.parameters == nil {
		self.parameters = map[string]interface{}{}
	}
	self.parameters[key] = value
	return self
}
func (self DynamoPut) GatherStates() []State {
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

type dynamoputOutput struct {
	Catch      []*CatchClause         `json:"Catch,omitempty"`
	Comment    string                 `json:"Comment,omitempty"`
	End        bool                   `json:"End,omitempty"`
	Heartbeat  Timeout                `json:"HeartbeatSeconds,omitempty"`
	Next       string                 `json:"Next,omitempty"`
	Parameters map[string]interface{} `json:"Parameters,omitempty"`
	Resource   string                 `json:"Resource,omitempty"`
	ResultPath string                 `json:"ResultPath,omitempty"`
	Timeout    Timeout                `json:"TimeoutSeconds,omitempty"`
	Type       StateType              `json:"Type,omitempty"`
}

func (self DynamoPut) StateType() StateType {
	return StateType("Task")
}
