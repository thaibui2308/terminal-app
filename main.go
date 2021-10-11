package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thaibui2308/terminal-app/api"
	"github.com/thaibui2308/terminal-app/models"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#D3D5D4"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF9ACA"))
	detailStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF9ACA"))
	blurredDetailStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8E4162"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EDA2C0"))

	focusedButton = focusedStyle.Copy().Render("[ 🙊Search ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("🙊Search"))

	detailButton        = detailStyle.Copy().Render("[ 👓More Details ]")
	blurredDetailButton = fmt.Sprintf("[ %s ]", blurredDetailStyle.Render("👓More Details"))

	docStyle = lipgloss.NewStyle().Margin(1, 2)

	studentButton = focusedStyle.Copy().Render("[ Ratings ]")
)

type model struct {
	// Input field
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode

	// Show the result of the search
	professor models.Professors
	list      list.Model

	// Details about a professor
	ratings []models.Ratings
	content string
	ready   bool

	notFound bool
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.NewModel()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "First Name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Last Name"
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}
func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.list.Items()) == 0 {
		if len(m.inputs) == 1 {
			return m, tea.Quit
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit

			// Change cursor mode
			case "ctrl+r":
				m.cursorMode++
				if m.cursorMode > textinput.CursorHide {
					m.cursorMode = textinput.CursorBlink
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := range m.inputs {
					cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
				}
				return m, tea.Batch(cmds...)

			// Set focus to next input
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				// Did the user press enter while the submit button was focused?
				// If so, search the professor.
				if s == "enter" && m.focusIndex == len(m.inputs) {
					m.focusIndex = 0
					firstName := m.inputs[0].Value()
					lastName := m.inputs[1].Value()
					m.updateSearchResults(firstName, lastName)
					return m, nil

				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				if m.focusIndex > len(m.inputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs)
				}

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}

				return m, tea.Batch(cmds...)
			}
		}

		// Handle character input and blinking
		cmd := m.updateInputs(msg)

		return m, cmd
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" {
				return m, nil
			} else if msg.String() == "tab" || msg.String() == "shift+tab" || msg.String() == "enter" || msg.String() == "up" || msg.String() == "down" {
				s := msg.String()

				// Did the user press enter while the submit button was focused?
				// If so, search the professor.
				if s == "enter" && m.focusIndex == len(m.list.Items()) {
					m.RequestForDetail(strconv.Itoa(m.professor.Tid))
					return m, nil

				}

				// Cycle indexes
				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				} else {
					m.focusIndex++
				}

				if m.focusIndex > len(m.list.Items()) {
					m.focusIndex = len(m.list.Items())
				} else if m.focusIndex < 0 {
					m.focusIndex = 0
				}

			}
		case tea.WindowSizeMsg:
			top, right, bottom, left := docStyle.GetMargin()
			m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		}
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder
	if len(m.list.Items()) == 0 {
		if m.notFound {
			b.WriteString(m.inputs[0].View())
			b.WriteString(helpStyle.Render("\nCannot find this instructor in the system!"))
			b.WriteString(helpStyle.Render("\n\nPress enter to exit"))
		} else {
			for i := range m.inputs {
				b.WriteString(m.inputs[i].View())
				if i < len(m.inputs)-1 {
					b.WriteRune('\n')
				}
			}

			button := &blurredButton
			if m.focusIndex == len(m.inputs) {
				button = &focusedButton
			}

			fmt.Fprintf(&b, "\n\n%s\n\n", *button)

			b.WriteString(helpStyle.Render("cursor mode is "))
			b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
			b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
			b.WriteString(helpStyle.Render("\nPress enter to search!"))

		}
		return b.String()
	} else {
		if len(m.ratings) == 0 {
			b.WriteString("\n" + m.list.View())
			button := &blurredDetailButton
			if m.focusIndex == 5 {
				button = &detailButton
			}
			fmt.Fprintf(&b, "\n%s\n\n", *button)
			b.WriteString(helpStyle.Render("\nPress enter to see more details!"))

		} else {
			button := &studentButton
			fmt.Fprintf(&b, "\n%s\n\n", *button)
			for _, v := range m.ratings {
				comment := strings.Split(v.RComments, " ")
				for i, word := range comment {
					if i == 0 {
						b.WriteString("🍓 " + word + " ")
					} else if (i+1)%15 == 0 {
						b.WriteString(word + "\n   ")
					} else {
						b.WriteString(word + " ")
					}
				}
				b.WriteString("\n")
			}
			b.WriteString(helpStyle.Render("Press q to exit"))
		}
		return b.String()
	}

}

func main() {
	Model := initialModel()
	p := tea.NewProgram(
		Model,
	)
	if err := p.Start(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

func (m *model) updateSearchResults(firstName, lastName string) {
	var professor models.Professors
	professor, err := api.FindInstructor(firstName, lastName)

	if err != nil {
		input := make([]textinput.Model, 1)
		var t textinput.Model
		for i := range input {
			t = textinput.NewModel()
			t.CursorStyle = cursorStyle
			t.CharLimit = 32
			switch i {
			case 0:
				t.Placeholder = ""
				t.Focus()
				t.PromptStyle = focusedStyle
				t.TextStyle = focusedStyle
			}
			input[i] = t
		}
		m.NotFound()
		m.inputs = input
	} else {
		items := []list.Item{
			item{title: "Institution name", desc: professor.InstitutionName},
			item{title: "Department", desc: professor.TDept},
			item{title: "Number of ratings", desc: strconv.Itoa(professor.TNumRatings)},
			item{title: "Overall quality", desc: professor.RatingClass},
			item{title: "Overall rating", desc: professor.OverallRating},
		}
		m.professor = professor
		m.list = list.NewModel(items, list.NewDefaultDelegate(), 0, 0)
		m.list.Title = "✔️FOUND!"
	}
}

func (m *model) RequestForDetail(tid string) {
	ratingList := api.GetRatings(tid)
	var content string
	for _, rating := range ratingList {
		content += rating.RComments + "\n\n"
	}
	m.ratings = ratingList
	m.content = content
	m.ready = true
}

func (m *model) NotFound() {
	m.notFound = true
}
