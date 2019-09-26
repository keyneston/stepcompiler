package stepcompiler

import "encoding/json"

const (
	StatesAll = "States.ALL"
)

type CatchClause struct {
	errorEquals []string
	resultPath  string
	next        State
}

type catchClaseOuput struct {
	ErrorEquals []string `json:"ErrorEquals,omitempty"`
	ResultPath  string   `json:"ResultPath,omitempty"`
	Next        string   `json:"Next,omitempty"`
	End         bool     `json:"End,omitempty"`
}

func NewCatchClause() *CatchClause {
	return &CatchClause{
		errorEquals: []string{StatesAll},
	}
}

func (cc *CatchClause) Next(state State) *CatchClause {
	cc.next = state

	return cc
}

func (cc *CatchClause) NextChainable(state State) {
	cc.Next(state)
}

func (cc *CatchClause) GatherStates() []State {
	res := []State{}

	if cc.next != nil {
		res = append(res, cc.next)
		res = append(res, cc.next.GatherStates()...)
	}

	return res
}

func (cc CatchClause) MarshalJSON() ([]byte, error) {
	out := catchClaseOuput{
		ErrorEquals: cc.errorEquals,
		ResultPath:  cc.resultPath,
	}

	if cc.next != nil {
		out.Next = cc.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}
