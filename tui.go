package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	styleField = iota
	layersField
	trunkLengthField
	angleField
	leafLengthField
	branchCharsField
	leafCharsField
	widthField
	heightField
	renderModeField
	waitTimeField
	fixedWindowField
	startField
)

var treeStyles = []string{
	"Random",
	"Classic",
	"Fibonacci",
	"Offset Fibonacci",
	"Random Offset Fibonacci",
}

var (
	tuiTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))
	tuiSubtitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))
	tuiCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("61")).
			Padding(0, 1)
	tuiSectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))
	tuiLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
	tuiValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("121"))
	tuiSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("61")).
				Bold(true)
	tuiHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
	tuiErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204"))
	tuiTreeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("114"))
)

type optionsTUIModel struct {
	options  *Options
	defaults *Options
	cursor   int
	editing  bool
	input    string
	error    string
	accepted bool
	width    int
	height   int
}

func newOptionsTUIModel(options *Options) optionsTUIModel {
	return optionsTUIModel{
		options:  cloneOptions(options),
		defaults: cloneOptions(options),
	}
}

func cloneOptions(options *Options) *Options {
	copy := *options
	return &copy
}

func runOptionsTUI(options *Options) (*Options, bool, error) {
	program := tea.NewProgram(newOptionsTUIModel(options), tea.WithAltScreen())
	result, err := program.Run()
	if err != nil {
		return nil, false, err
	}

	model, ok := result.(optionsTUIModel)
	if !ok {
		return nil, false, fmt.Errorf("unexpected options interface result")
	}
	return model.options, model.accepted, nil
}

func isInteractiveTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

func shouldUseOptionsTUI(args []string) bool {
	return hasTUIArg(args) || (len(args) == 0 && isInteractiveTerminal())
}

func hasTUIArg(args []string) bool {
	for _, arg := range args {
		if arg == "--tui" {
			return true
		}
	}
	return false
}

func withoutTUIArg(args []string) []string {
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "--tui" {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}

func (m optionsTUIModel) Init() tea.Cmd {
	return nil
}

func (m optionsTUIModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.editing {
			return m.updateTextInput(msg)
		}
		return m.updateSelection(msg)
	}

	return m, nil
}

func (m optionsTUIModel) updateTextInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if len([]rune(m.input)) == 0 {
			m.error = "Choose at least one character."
			return m, nil
		}
		m.setTextValue(m.input)
		m.editing = false
		m.error = ""
	case tea.KeyEsc:
		m.editing = false
		m.error = ""
	case tea.KeyBackspace:
		chars := []rune(m.input)
		if len(chars) > 0 {
			m.input = string(chars[:len(chars)-1])
		}
	case tea.KeyRunes:
		m.input += string(msg.Runes)
	}
	return m, nil
}

func (m optionsTUIModel) updateSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch key {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "up", "k", "shift+tab":
		m.cursor = (m.cursor - 1 + startField + 1) % (startField + 1)
	case "down", "j", "tab":
		m.cursor = (m.cursor + 1) % (startField + 1)
	case "left", "h", "-":
		m.adjust(-1)
	case "right", "l", "+", "=":
		m.adjust(1)
	case "r":
		m.options = cloneOptions(m.defaults)
		m.error = "Settings restored."
	case "enter", " ":
		if m.cursor == startField {
			m.accepted = true
			return m, tea.Quit
		}
		if m.cursor == branchCharsField || m.cursor == leafCharsField {
			m.editing = true
			m.input = m.textValue()
			m.error = ""
		} else {
			m.adjust(1)
		}
	}

	return m, nil
}

func (m *optionsTUIModel) adjust(delta int) {
	m.error = ""
	switch m.cursor {
	case styleField:
		selection := (m.treeStyleSelection() + delta + len(treeStyles)) % len(treeStyles)
		m.setTreeStyle(selection)
	case layersField:
		m.options.NumLayers = clamp(m.options.NumLayers+delta, 1, 12)
	case trunkLengthField:
		m.options.InitialLen = clamp(m.options.InitialLen+delta, 1, 30)
	case angleField:
		degrees := int(math.Round(m.options.AngleMean * 180 / math.Pi))
		m.options.AngleMean = float64(clamp(degrees+delta, 5, 85)) * math.Pi / 180
	case leafLengthField:
		m.options.LeafLen = clamp(m.options.LeafLen+delta, 1, 10)
	case widthField:
		m.options.WindowWidth = clamp(m.options.WindowWidth+delta*2, 30, 160)
	case heightField:
		m.options.WindowHeight = clamp(m.options.WindowHeight+delta, 12, 60)
	case renderModeField:
		m.options.Instant = !m.options.Instant
	case waitTimeField:
		m.options.WaitTime = math.Round(clampFloat(m.options.WaitTime+float64(delta)*0.01, 0, 1)*100) / 100
	case fixedWindowField:
		m.options.FixedWindow = !m.options.FixedWindow
	}
}

func clamp(value, lower, upper int) int {
	return min(max(value, lower), upper)
}

func clampFloat(value, lower, upper float64) float64 {
	return math.Min(math.Max(value, lower), upper)
}

func (m *optionsTUIModel) setTreeStyle(selection int) {
	if selection == 0 {
		m.options.UserSetType = false
		m.options.Type = rng.IntN(4)
		return
	}
	m.options.UserSetType = true
	m.options.Type = selection - 1
}

func (m optionsTUIModel) treeStyleSelection() int {
	if !m.options.UserSetType {
		return 0
	}
	return clamp(m.options.Type+1, 1, len(treeStyles)-1)
}

func (m optionsTUIModel) textValue() string {
	if m.cursor == branchCharsField {
		return m.options.BranchChars
	}
	return m.options.LeafChars
}

func (m *optionsTUIModel) setTextValue(value string) {
	if m.cursor == branchCharsField {
		m.options.BranchChars = value
		return
	}
	m.options.LeafChars = value
}

func (m optionsTUIModel) fieldValue(field int) string {
	switch field {
	case styleField:
		return treeStyles[m.treeStyleSelection()]
	case layersField:
		return fmt.Sprintf("%d", m.options.NumLayers)
	case trunkLengthField:
		return fmt.Sprintf("%d", m.options.InitialLen)
	case angleField:
		return fmt.Sprintf("%d°", int(math.Round(m.options.AngleMean*180/math.Pi)))
	case leafLengthField:
		return fmt.Sprintf("%d", m.options.LeafLen)
	case branchCharsField:
		return fmt.Sprintf("%q", m.options.BranchChars)
	case leafCharsField:
		return fmt.Sprintf("%q", m.options.LeafChars)
	case widthField:
		return fmt.Sprintf("%d columns", m.options.WindowWidth)
	case heightField:
		return fmt.Sprintf("%d rows", m.options.WindowHeight)
	case renderModeField:
		if m.options.Instant {
			return "Instant"
		}
		return "Animated"
	case waitTimeField:
		return fmt.Sprintf("%.2f seconds", m.options.WaitTime)
	case fixedWindowField:
		if m.options.FixedWindow {
			return "Keep height fixed"
		}
		return "Expand when needed"
	case startField:
		return "Grow bonsai  →"
	default:
		return ""
	}
}

func (m optionsTUIModel) fieldRow(field int, label string) string {
	value := m.fieldValue(field)
	if m.editing && m.cursor == field {
		value = m.input + "▏"
	}

	row := tuiLabelStyle.Width(17).Render(label) + tuiValueStyle.Render(value)
	if m.cursor == field {
		return tuiSelectedStyle.Width(43).Render("› " + row)
	}
	return "  " + row
}

func (m optionsTUIModel) settingsCard(title string, fields []struct {
	field int
	label string
}) string {
	rows := make([]string, 0, len(fields)+1)
	rows = append(rows, tuiSectionStyle.Render(title))
	for _, field := range fields {
		rows = append(rows, m.fieldRow(field.field, field.label))
	}
	return tuiCardStyle.Render(strings.Join(rows, "\n"))
}

func (m optionsTUIModel) compactSettingsCard() string {
	sections := []struct {
		title  string
		fields []struct {
			field int
			label string
		}
	}{
		{
			title: "TREE",
			fields: []struct {
				field int
				label string
			}{
				{styleField, "Style"},
				{layersField, "Branch layers"},
				{trunkLengthField, "Trunk length"},
				{angleField, "Branch angle"},
				{leafLengthField, "Leaf size"},
			},
		},
		{
			title: "CHARACTER",
			fields: []struct {
				field int
				label string
			}{
				{branchCharsField, "Branch glyphs"},
				{leafCharsField, "Leaf glyphs"},
			},
		},
		{
			title: "CANVAS & MOTION",
			fields: []struct {
				field int
				label string
			}{
				{widthField, "Canvas width"},
				{heightField, "Canvas height"},
				{renderModeField, "Render mode"},
				{waitTimeField, "Draw delay"},
				{fixedWindowField, "Overflow"},
				{startField, ""},
			},
		},
	}

	rows := make([]string, 0, startField+4)
	for _, section := range sections {
		rows = append(rows, tuiSectionStyle.Render(section.title))
		for _, field := range section.fields {
			rows = append(rows, m.fieldRow(field.field, field.label))
		}
	}
	return tuiCardStyle.Render(strings.Join(rows, "\n"))
}

func (m optionsTUIModel) View() string {
	if m.width > 0 && m.width < 55 {
		return "Your terminal is too narrow for the bonsai configurator. Please use at least 55 columns.\n"
	}

	shape := []string{
		"       .-&-.       ",
		"     .- /|\\ -.     ",
		"       / | \\       ",
		"      /  |  \\      ",
		"    .____|____.     ",
		"     \\  bonsai /     ",
	}

	fullHeader := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tuiTreeStyle.Render(strings.Join(shape, "\n")),
		"  ",
		lipgloss.JoinVertical(
			lipgloss.Left,
			tuiTitleStyle.Render("GoBonsai · Shape a small world"),
			tuiSubtitleStyle.Render("Tune your tree, then let it grow."),
		),
	)

	compact := m.width < 110 || (m.height > 0 && m.height < 32)
	if compact {
		header := tuiTitleStyle.Render("GoBonsai") + tuiSubtitleStyle.Render("  ·  Shape a small world")
		status := "↑/↓ navigate  •  ←/→ adjust  •  enter edit/select  •  r reset  •  q cancel"
		if m.editing {
			status = "Type characters  •  enter save  •  esc discard"
		}
		if m.error != "" {
			if m.error == "Settings restored." {
				status = tuiHintStyle.Render(m.error)
			} else {
				status = tuiErrorStyle.Render(m.error)
			}
		} else {
			status = tuiHintStyle.Render(status)
		}

		return "\n" + header + "\n\n" + m.compactSettingsCard() + "\n\n" + status + "\n"
	}

	treeCard := m.settingsCard("TREE", []struct {
		field int
		label string
	}{
		{styleField, "Style"},
		{layersField, "Branch layers"},
		{trunkLengthField, "Trunk length"},
		{angleField, "Branch angle"},
		{leafLengthField, "Leaf size"},
	})

	lookCard := m.settingsCard("CHARACTER", []struct {
		field int
		label string
	}{
		{branchCharsField, "Branch glyphs"},
		{leafCharsField, "Leaf glyphs"},
	})

	canvasCard := m.settingsCard("CANVAS & MOTION", []struct {
		field int
		label string
	}{
		{widthField, "Canvas width"},
		{heightField, "Canvas height"},
		{renderModeField, "Render mode"},
		{waitTimeField, "Draw delay"},
		{fixedWindowField, "Overflow"},
		{startField, ""},
	})

	var cards string
	left := lipgloss.JoinVertical(lipgloss.Left, treeCard, lookCard)
	cards = lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", canvasCard)

	status := "↑/↓ navigate  •  ←/→ adjust  •  enter edit/select  •  r reset  •  q cancel"
	if m.editing {
		status = "Type characters  •  enter save  •  esc discard"
	}
	if m.error != "" {
		status = m.error
		if m.error == "Settings restored." {
			status = tuiHintStyle.Render(status)
		} else {
			status = tuiErrorStyle.Render(status)
		}
	} else {
		status = tuiHintStyle.Render(status)
	}

	return "\n" + fullHeader + "\n\n" + cards + "\n\n" + status + "\n"
}
