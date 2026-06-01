package input

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func key(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func specialKey(t tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: t}
}

func TestDefaultKeymap_BasicActions(t *testing.T) {
	km := DefaultKeymap

	tests := []struct {
		input  tea.KeyMsg
		action Action
	}{
		{key("j"), ActionDown},
		{key("k"), ActionUp},
		{key("q"), ActionQuit},
		{key("/"), ActionSearch},
		{key("n"), ActionNextHit},
		{key("N"), ActionPrevHit},
		{key("r"), ActionReload},
		{key("?"), ActionHelp},
		{key("G"), ActionBottom},
		{specialKey(tea.KeyEsc), ActionEsc},
		{specialKey(tea.KeyEnter), ActionEnter},
	}

	for _, tt := range tests {
		action, _ := km.Resolve(tt.input)
		if action != tt.action {
			t.Errorf("key %q: want action %d, got %d", tt.input.String(), tt.action, action)
		}
	}
}

func TestDefaultKeymap_HalfPage(t *testing.T) {
	km := DefaultKeymap

	upMsg := tea.KeyMsg{Type: tea.KeyCtrlU}
	downMsg := tea.KeyMsg{Type: tea.KeyCtrlD}

	action, _ := km.Resolve(upMsg)
	if action != ActionHalfUp {
		t.Errorf("ctrl+u: want ActionHalfUp, got %d", action)
	}

	action, _ = km.Resolve(downMsg)
	if action != ActionHalfDown {
		t.Errorf("ctrl+d: want ActionHalfDown, got %d", action)
	}
}

func TestDefaultKeymap_UnknownKey(t *testing.T) {
	km := DefaultKeymap
	action, ch := km.Resolve(key("z"))
	if action != ActionChar {
		t.Errorf("want ActionChar, got %d", action)
	}
	if ch != "z" {
		t.Errorf("want ch='z', got %q", ch)
	}
}

func TestDefaultKeymap_CaseInsensitive(t *testing.T) {
	km := DefaultKeymap
	
	action, _ := km.Resolve(key("J"))
	if action != ActionDown {
		t.Errorf("'J' should resolve to ActionDown (case-insensitive), got %d", action)
	}
}

func TestGGBuffer_SinglePress(t *testing.T) {
	var g GGBuffer
	
	if g.Feed("g", "g") {
		t.Error("single 'g' should not trigger top")
	}
}

func TestGGBuffer_DoublePress(t *testing.T) {
	var g GGBuffer
	g.Feed("g", "g")
	
	if !g.Feed("g", "g") {
		t.Error("double 'g' within timeout should trigger top")
	}
}

func TestGGBuffer_Expired(t *testing.T) {
	var g GGBuffer
	g.Feed("g", "g")
	
	g.deadline = time.Now().Add(-time.Second)
	
	if g.Feed("g", "g") {
		t.Error("double 'g' after timeout should NOT trigger top")
	}
}

func TestGGBuffer_WrongKey(t *testing.T) {
	var g GGBuffer
	g.Feed("g", "g")
	
	if g.Feed("x", "g") {
		t.Error("different key should not trigger top")
	}
	if g.pending {
		t.Error("pending should be reset after wrong key")
	}
}

func TestUIState_Values(t *testing.T) {
	if StateNormal != 0 {
		t.Errorf("StateNormal should be 0, got %d", StateNormal)
	}
	if StateSearch == StateNormal {
		t.Error("StateSearch should differ from StateNormal")
	}
	if StateHelp == StateNormal || StateHelp == StateSearch {
		t.Error("StateHelp should be unique")
	}
}

func TestKeymap_CustomKeys(t *testing.T) {
	km := Keymap{
		Up:       "w",
		Down:     "s",
		HalfUp:   "ctrl+u",
		HalfDown: "ctrl+d",
		Top:      "g",
		Bottom:   "G",
		Search:   "/",
		NextHit:  "n",
		PrevHit:  "N",
		Reload:   "r",
		Quit:     "x",
		Help:     "h",
	}

	action, _ := km.Resolve(key("w"))
	if action != ActionUp {
		t.Errorf("custom up key 'w': want ActionUp, got %d", action)
	}

	action, _ = km.Resolve(key("s"))
	if action != ActionDown {
		t.Errorf("custom down key 's': want ActionDown, got %d", action)
	}

	action, _ = km.Resolve(key("x"))
	if action != ActionQuit {
		t.Errorf("custom quit key 'x': want ActionQuit, got %d", action)
	}
}