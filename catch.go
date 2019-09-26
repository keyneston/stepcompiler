package stepcompiler

type CatchClause struct {
	ErrorEquals []string `json:"ErrorEquals"`
	ResultPath  string   `json:"ResultPath"`
	Next        string   `json:"Next"`
}
