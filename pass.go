package stepcompiler

import "encoding/json"

type Pass struct {
	name       string
	result     map[string]interface{}
	resultPath string
	parameters map[string]interface{}

	comment string
	next    State
}

type passOutput struct {
	Type       StateType              `json:"Type"`
	Parameters map[string]interface{} `json:"Parameters,omitempty"`
	ResultPath string                 `json:"ResultPath,omitempty"`
	Result     map[string]interface{} `json:"Result,omitempty"`
	Next       string                 `json:"Next,omitempty"`
	Comment    string                 `json:"Comment,omitempty"`
	End        bool                   `json:"End,omitempty"`
}

func NewPass(name string) *Pass {
	return &Pass{
		name: name,
	}
}

func (p *Pass) Next(state State) *Pass {
	p.next = state

	return p
}

func (p *Pass) NextChainable(state State) {
	p.Next(state)
}

func (p *Pass) ResultPath(resultPath string) *Pass {
	p.resultPath = resultPath
	return p
}

func (p *Pass) GatherStates() []State {
	res := []State{}

	if p.next != nil {
		res = append(res, p.next)
	}

	return res
}

func (Pass) StateType() StateType {
	return PassType
}

func (p *Pass) Result(input map[string]interface{}) *Pass {
	p.result = input
	return p
}

func (p Pass) Name() string {
	return p.name
}

func (p Pass) MarshalJSON() ([]byte, error) {
	out := passOutput{
		Type:       p.StateType(),
		Comment:    p.comment,
		Result:     p.result,
		ResultPath: p.resultPath,
	}

	if p.next != nil {
		out.Next = p.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}
