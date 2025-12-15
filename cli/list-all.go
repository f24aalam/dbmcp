package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/f24aalam/godbmcp/storage"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list list.Model
	selected *storage.Credential
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "ctr+c":
					return m, tea.Quit
				case "enter":
					selected := m.list.SelectedItem()
					if selected != nil {
						cred := selected.(storage.Credential)
						m.selected = &cred

						return m, tea.Quit
					}
			}
		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width - h, msg.Height - v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func ListAllConnections() {
	creds, err := storage.ListCredentials()

	if err != nil {
		fmt.Println("Error fetching: ", err)
		return
	}

	items := []list.Item{}
	for _, cred := range creds {
		items = append(items, list.Item(cred))
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Stored Connections"

	p := tea.NewProgram(m, tea.WithAltScreen())
	updatedModel, err := p.Run()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}

	m = updatedModel.(model)
	if m.selected != nil {
		ShowOptions(m.selected)
	}
}

func ShowOptions(cred *storage.Credential) {
	var option string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Options for %s", cred.ID)).
				Options(
					huh.NewOption("Edit", "edit"),
					huh.NewOption("Delete", "delete"),
					huh.NewOption("Connect", "connect"),
				).
				Value(&option),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("Error rendering options: ", err)
		return
	}

	switch option {
	case "edit":
		AddNewConnection(cred)
	}
}
