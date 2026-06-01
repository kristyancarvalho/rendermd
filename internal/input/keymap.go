package input

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Keymap struct {
	Up       string
	Down     string
	HalfUp   string
	HalfDown string
	Top      string
	Bottom   string
	Search   string
	NextHit  string
	PrevHit  string
	Reload   string
	Quit     string
	Help     string
}

var DefaultKeymap = Keymap{
	Up:       "k",
	Down:     "j",
	HalfUp:   "ctrl+u",
	HalfDown: "ctrl+d",
	Top:      "g",
	Bottom:   "G",
	Search:   "/",
	NextHit:  "n",
	PrevHit:  "N",
	Reload:   "r",
	Quit:     "q",
	Help:     "?",
}

type Action int

const (
	ActionNone Action = iota
	ActionUp
	ActionDown
	ActionHalfUp
	ActionHalfDown
	ActionTop
	ActionBottom
	ActionSearch
	ActionNextHit
	ActionPrevHit
	ActionReload
	ActionQuit
	ActionHelp
	ActionEsc
	ActionEnter
	ActionChar
)

type GGBuffer struct {
	pending  bool
	deadline time.Time
}

func (g *GGBuffer) Feed(ch string, topKey string) bool {
	if ch != topKey {
		g.pending = false
		return false
	}
	if g.pending && time.Now().Before(g.deadline) {
		g.pending = false
		return true
	}
	g.pending = true
	g.deadline = time.Now().Add(500 * time.Millisecond)
	return false
}

func (k Keymap) Resolve(msg tea.KeyMsg) (Action, string) {
	key := msg.String()
	switch {
	case key == "esc":
		return ActionEsc, ""
	case key == "enter":
		return ActionEnter, ""
	case matchKey(key, k.Up):
		return ActionUp, ""
	case matchKey(key, k.Down):
		return ActionDown, ""
	case matchKey(key, k.HalfUp):
		return ActionHalfUp, ""
	case matchKey(key, k.HalfDown):
		return ActionHalfDown, ""
	case matchKey(key, k.Bottom):
		return ActionBottom, ""
	case matchKey(key, k.Search):
		return ActionSearch, ""
	case matchKey(key, k.PrevHit):
		return ActionPrevHit, ""
	case matchKey(key, k.NextHit):
		return ActionNextHit, ""
	case matchKey(key, k.Reload):
		return ActionReload, ""
	case matchKey(key, k.Quit):
		return ActionQuit, ""
	case matchKey(key, k.Help):
		return ActionHelp, ""
	}
	return ActionChar, key
}

func matchKey(got, want string) bool {
	return strings.EqualFold(got, want)
}