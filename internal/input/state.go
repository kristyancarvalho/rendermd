package input

type UIState int

const (
	StateNormal UIState = iota
	StateSearch
	StatePrompt
	StateHelp
)
