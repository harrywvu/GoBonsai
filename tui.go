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
	branchCharsField
	leafCharsField
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
			Foreground(lipgloss.Color("151"))
	tuiSubtitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244"))
	tuiCardStyle = lipgloss.NewStyle().
			Width(50).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("65")).
			Padding(1, 1)
	tuiSectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("180"))
	tuiLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
	tuiValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("151"))
	tuiSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("65")).
				Bold(true).
				Width(50)
	tuiGrowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("65")).
			Bold(true).
			Width(50).
			Align(lipgloss.Center)
	tuiHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
	tuiErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("210"))
	tuiTreeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("107"))
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
	switch msg.String() {
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
	}
}

func clamp(value, lower, upper int) int {
	return min(max(value, lower), upper)
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
	case branchCharsField:
		return fmt.Sprintf("%q", m.options.BranchChars)
	case leafCharsField:
		return fmt.Sprintf("%q", m.options.LeafChars)
	default:
		return ""
	}
}

func (m optionsTUIModel) fieldRow(field int, label string) string {
	value := m.fieldValue(field)
	if m.editing && m.cursor == field {
		value = m.input + "▏"
	}

	row := tuiLabelStyle.Width(19).Render(label) + tuiValueStyle.Render(value)
	if m.cursor == field {
		return tuiSelectedStyle.Render("› " + row)
	}
	return "  " + row
}

func (m optionsTUIModel) View() string {
	if m.width > 0 && m.width < 58 {
		return "Your terminal is too narrow for the bonsai configurator. Please use at least 58 columns.\n"
	}

	bonsai := tuiTreeStyle.Render(strings.Join([]string{
		"    .-.-.  ",
		"  .(  |  ).",
		"     / \\   ",
		"   .-===-. ",
	}, "\n"))
	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		bonsai,
		"  ",
		lipgloss.JoinVertical(
			lipgloss.Left,
			tuiTitleStyle.Render("GoBonsai"),
			tuiSubtitleStyle.Render("Give your tree a little character."),
		),
	)

	rows := []string{
		header,
		"",
		tuiSectionStyle.Render("TREE SHAPE"),
		m.fieldRow(styleField, "Tree style"),
		m.fieldRow(layersField, "Branch layers"),
		m.fieldRow(trunkLengthField, "Trunk length"),
		m.fieldRow(angleField, "Branch angle"),
		"",
		tuiSectionStyle.Render("TEXTURE"),
		m.fieldRow(branchCharsField, "Branch glyphs"),
		m.fieldRow(leafCharsField, "Leaf glyphs"),
		"",
	}
	if m.cursor == startField {
		rows = append(rows, tuiGrowStyle.Render("✦  Grow bonsai  ✦"))
	} else {
		rows = append(rows, tuiHintStyle.Width(50).Align(lipgloss.Center).Render("Grow bonsai"))
	}

	status := "↑↓ move · ←→ tune · enter · r reset · q quit"
	if m.editing {
		status = "Type glyphs · enter save · esc discard"
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
	rows = append(rows, "", status)

	dialog := tuiCardStyle.Render(strings.Join(rows, "\n"))
	if m.width == 0 || m.height == 0 {
		return "\n" + dialog + "\n"
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialog)
}
