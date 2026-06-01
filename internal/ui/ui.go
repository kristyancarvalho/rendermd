package ui

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kristyancarvalho/mdp/internal/config"
	"github.com/kristyancarvalho/mdp/internal/input"
	"github.com/kristyancarvalho/mdp/internal/layout"
	"github.com/kristyancarvalho/mdp/internal/markdown"
	"github.com/kristyancarvalho/mdp/internal/model"
	"github.com/kristyancarvalho/mdp/internal/render"
	"github.com/kristyancarvalho/mdp/internal/theme"
	"github.com/kristyancarvalho/mdp/internal/watch"
)

type watchMsg struct{}
type errMsg struct{ err error }

type uiModel struct {
	filename  string
	content   []byte
	doc       model.Document
	lines     []layout.Line
	lastValid []layout.Line
	engine    *layout.Engine
	viewport  render.Viewport
	renderer  *render.Renderer
	state     input.UIState
	keymap    input.Keymap
	ggBuf     input.GGBuffer
	query     string
	queryBuf  string
	hits      []int
	hitIndex  int
	thm       theme.Theme
	cfg       config.Config
	watcher   *watch.Watcher
	lastErr   error
}

func Run(filename string, content []byte, cfg config.Config, w *watch.Watcher) error {
	thm := cfg.ResolvedTheme()
	doc := markdown.Parse(content)
	eng := &layout.Engine{}

	m := &uiModel{
		filename: filename,
		content:  content,
		doc:      doc,
		engine:   eng,
		renderer: render.New(thm),
		state:    input.StateNormal,
		keymap:   toKeymap(cfg.Keys),
		thm:      thm,
		cfg:      cfg,
		watcher:  w,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if w != nil {
		go func() {
			for {
				select {
				case <-w.Events:
					p.Send(watchMsg{})
				case err := <-w.Errors:
					p.Send(errMsg{err: err})
				}
			}
		}()
	}

	_, err := p.Run()
	return err
}

func (m *uiModel) Init() tea.Cmd {
	return nil
}

func (m *uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 1
		m.reLayout()
	case watchMsg:
		data, err := os.ReadFile(m.filename)
		if err != nil {
			m.lastErr = err
		} else {
			m.content = data
			m.doc = markdown.Parse(data)
			m.reLayout()
			m.lastErr = nil
			m.lastValid = m.lines
		}
	case errMsg:
		m.lastErr = msg.err
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m *uiModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == input.StateSearch {
		return m.handleSearchKey(msg)
	}
	if m.state == input.StateHelp {
		action, _ := m.keymap.Resolve(msg)
		if action == input.ActionQuit || action == input.ActionHelp || action == input.ActionEsc {
			m.state = input.StateNormal
		}
		return m, nil
	}

	action, ch := m.keymap.Resolve(msg)
	switch action {
	case input.ActionQuit:
		return m, tea.Quit
	case input.ActionDown:
		m.scrollBy(1)
	case input.ActionUp:
		m.scrollBy(-1)
	case input.ActionHalfDown:
		m.scrollBy(m.viewport.Height / 2)
	case input.ActionHalfUp:
		m.scrollBy(-(m.viewport.Height / 2))
	case input.ActionBottom:
		m.viewport.Offset = len(m.lines) - m.viewport.Height
		if m.viewport.Offset < 0 {
			m.viewport.Offset = 0
		}
	case input.ActionSearch:
		m.state = input.StateSearch
		m.queryBuf = ""
	case input.ActionNextHit:
		m.nextHit(1)
	case input.ActionPrevHit:
		m.nextHit(-1)
	case input.ActionReload:
		if m.filename != "" {
			data, err := os.ReadFile(m.filename)
			if err != nil {
				m.lastErr = err
			} else {
				m.content = data
				m.doc = markdown.Parse(data)
				m.engine = &layout.Engine{}
				m.reLayout()
				m.lastErr = nil
			}
		}
	case input.ActionHelp:
		m.state = input.StateHelp
	case input.ActionChar:
		if m.ggBuf.Feed(ch, m.keymap.Top) {
			m.viewport.Offset = 0
		}
	}
	return m, nil
}

func (m *uiModel) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "esc":
		m.state = input.StateNormal
		m.queryBuf = ""
		m.query = ""
		m.hits = nil
	case "enter":
		m.query = m.queryBuf
		m.hits = findHits(m.lines, m.query)
		m.hitIndex = 0
		m.state = input.StateNormal
		if len(m.hits) > 0 {
			m.viewport.Offset = m.hits[0]
		}
	case "backspace":
		if len(m.queryBuf) > 0 {
			runes := []rune(m.queryBuf)
			m.queryBuf = string(runes[:len(runes)-1])
		}
	default:
		if utf8.RuneCountInString(key) == 1 {
			m.queryBuf += key
		}
	}
	return m, nil
}

func (m *uiModel) View() string {
	if m.viewport.Height <= 0 {
		return ""
	}
	if m.state == input.StateHelp {
		return m.helpView()
	}
	body := m.renderer.Render(m.lines, m.viewport)
	bar := m.statusBar()
	return body + "\n" + bar
}

func (m *uiModel) statusBar() string {
	thm := m.thm
	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(thm.Border)).
		Foreground(lipgloss.Color(thm.Text)).
		Width(m.viewport.Width)

	name := m.filename
	if name == "" {
		name = "stdin"
	}
	line := m.viewport.Offset + 1
	total := len(m.lines)
	left := fmt.Sprintf(" %s  [%d/%d]", name, line, total)

	var right string
	switch m.state {
	case input.StateSearch:
		right = fmt.Sprintf("/ %s_", m.queryBuf)
	default:
		if m.lastErr != nil {
			right = lipgloss.NewStyle().
				Foreground(lipgloss.Color(thm.Muted)).
				Render(m.lastErr.Error())
		} else if m.query != "" {
			if len(m.hits) == 0 {
				right = "[no results]"
			} else {
				right = fmt.Sprintf("[%d/%d]", m.hitIndex+1, len(m.hits))
			}
		}
	}

	gap := m.viewport.Width - utf8.RuneCountInString(left) - utf8.RuneCountInString(right)
	if gap < 1 {
		gap = 1
	}
	return barStyle.Render(left + strings.Repeat(" ", gap) + right)
}

func (m *uiModel) helpView() string {
	thm := m.thm
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(thm.Border)).
		Padding(1, 2).
		Foreground(lipgloss.Color(thm.Text))

	km := m.keymap
	content := fmt.Sprintf(
		"  %-12s scroll up\n"+
			"  %-12s scroll down\n"+
			"  %-12s half page up\n"+
			"  %-12s half page down\n"+
			"  %-12s top\n"+
			"  %-12s bottom\n"+
			"  %-12s search\n"+
			"  %-12s next result\n"+
			"  %-12s prev result\n"+
			"  %-12s reload\n"+
			"  %-12s help\n"+
			"  %-12s quit",
		km.Up, km.Down, km.HalfUp, km.HalfDown,
		km.Top+km.Top, km.Bottom,
		km.Search, km.NextHit, km.PrevHit,
		km.Reload, km.Help, km.Quit,
	)

	panel := panelStyle.Render(content)
	panelLines := strings.Split(panel, "\n")
	panelH := len(panelLines)
	panelW := 0
	for _, l := range panelLines {
		if utf8.RuneCountInString(l) > panelW {
			panelW = utf8.RuneCountInString(l)
		}
	}

	topPad := (m.viewport.Height - panelH) / 2
	leftPad := (m.viewport.Width - panelW) / 2
	if topPad < 0 {
		topPad = 0
	}
	if leftPad < 0 {
		leftPad = 0
	}

	var sb strings.Builder
	for i := 0; i < topPad; i++ {
		sb.WriteString("\n")
	}
	for _, l := range panelLines {
		sb.WriteString(strings.Repeat(" ", leftPad))
		sb.WriteString(l)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m *uiModel) reLayout() {
	cfg := layout.LayoutConfig{
		Width:      m.viewport.Width,
		Padding:    m.cfg.UI.Padding,
		SoftWrap:   m.cfg.UI.SoftWrap,
		ShowURLs:   m.cfg.UI.ShowURLs,
		HideSyntax: m.cfg.Markdown.HideSyntax,
	}
	m.lines = m.engine.Render(m.doc, cfg)
	if m.query != "" {
		m.hits = findHits(m.lines, m.query)
	}
}

func (m *uiModel) scrollBy(delta int) {
	m.viewport.Offset += delta
	maxOffset := len(m.lines) - m.viewport.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.viewport.Offset < 0 {
		m.viewport.Offset = 0
	}
	if m.viewport.Offset > maxOffset {
		m.viewport.Offset = maxOffset
	}
}

func (m *uiModel) nextHit(dir int) {
	if len(m.hits) == 0 {
		return
	}
	m.hitIndex = (m.hitIndex + dir + len(m.hits)) % len(m.hits)
	m.viewport.Offset = m.hits[m.hitIndex]
}

func findHits(lines []layout.Line, query string) []int {
	if query == "" {
		return nil
	}
	q := strings.ToLower(query)
	var hits []int
	for i, line := range lines {
		for _, seg := range line.Segments {
			if strings.Contains(strings.ToLower(seg.Text), q) {
				hits = append(hits, i)
				break
			}
		}
	}
	return hits
}

func toKeymap(k config.KeysConfig) input.Keymap {
	km := input.DefaultKeymap
	if k.Up != "" {
		km.Up = k.Up
	}
	if k.Down != "" {
		km.Down = k.Down
	}
	if k.HalfUp != "" {
		km.HalfUp = k.HalfUp
	}
	if k.HalfDown != "" {
		km.HalfDown = k.HalfDown
	}
	if k.Top != "" {
		km.Top = k.Top
	}
	if k.Bottom != "" {
		km.Bottom = k.Bottom
	}
	if k.Search != "" {
		km.Search = k.Search
	}
	if k.NextHit != "" {
		km.NextHit = k.NextHit
	}
	if k.PrevHit != "" {
		km.PrevHit = k.PrevHit
	}
	if k.Reload != "" {
		km.Reload = k.Reload
	}
	if k.Quit != "" {
		km.Quit = k.Quit
	}
	if k.Help != "" {
		km.Help = k.Help
	}
	return km
}

func clampMin(v, min int) int {
	if v < min {
		return min
	}
	return v
}