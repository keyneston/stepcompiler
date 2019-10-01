package stepcompiler

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDelete struct {
	name       string
	comment    string
	timeout    time.Duration
	next       State
	catch      []*CatchClause
	parameters map[string]interface{}
}

func NewDynamoDelete(name string) *DynamoDelete {
	return &DynamoDelete{
		name: name,
	}
}

func (DynamoDelete) StateType() StateType {
	return TaskType
}

func (dd *DynamoDelete) ConditionExpression(expression string) *DynamoDelete {
	return dd.setParameter("ConditionExpression", expression)
}

func (dd *DynamoDelete) TableName(tablename string) *DynamoDelete {
	return dd.setParameter("TableName", tablename)
}

func (dd *DynamoDelete) Key(values map[string]*dynamodb.AttributeValue) *DynamoDelete {
	return dd.setParameter("Key", values)
}

func (dd *DynamoDelete) setParameter(name string, value interface{}) *DynamoDelete {
	if dd.parameters == nil {
		dd.parameters = make(map[string]interface{})
	}

	dd.parameters[name] = value
	return dd
}

func (dd DynamoDelete) Name() string {
	return dd.name
}

func (dd *DynamoDelete) Comment(comment string) *DynamoDelete {
	dd.comment = comment
	return dd
}

func (dd *DynamoDelete) Next(state State) *DynamoDelete {
	dd.next = state
	return dd
}

func (dd *DynamoDelete) NextChainable(state State) {
	dd.Next(state)
}

func (dd *DynamoDelete) MarshalJSON() ([]byte, error) {
	out := taskOutput{
		Resource:       DynamoDeleteItemARN,
		Comment:        dd.comment,
		TimeoutSeconds: Timeout(dd.timeout),
		Type:           dd.StateType(),
		Parameters:     dd.parameters,
	}

	if dd.next != nil {
		out.Next = dd.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}

func (dd *DynamoDelete) GatherStates() []State {
	res := []State{dd}

	if dd.next != nil {
		res = append(res, dd.next)
		res = append(res, dd.next.GatherStates()...)
	}

	for _, clause := range dd.catch {
		res = append(res, clause.GatherStates()...)
	}

	return res
}

func (dd *DynamoDelete) Timeout(timeout time.Duration) *DynamoDelete {
	dd.timeout = timeout
	return dd
}
