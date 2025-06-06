package main

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type args struct {
	GeminiAPIKey string `arg:"-k,--api-key,required" help:"Google Gemini API key"`
	VaultRoot    string `arg:"positional,required" help:"root directory of the obsidian vault"`
}

func (args) Description() string {
	return "Feeds note to an LLM and create a note taking into account note content and linked websites."
}

func (args) Version() string {
	return "version 1"
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, description string
	Note               Note
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	selectItem key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		selectItem: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select item"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	selectedItem *item
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.selectItem) {
			selected := m.list.SelectedItem()
			if selected != nil {
				if itm, ok := selected.(item); ok {
					m.selectedItem = &itm
					return m, tea.Quit
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func newNoteList(items []list.Item, keys *listKeyMap) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "noto marked notes"

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.selectItem}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.selectItem}
	}

	return l
}

func run() error {
	var cfg args
	arg.MustParse(&cfg)

	scanner := NewScanner(cfg.VaultRoot)
	transformer := NewTransformer(cfg.VaultRoot, cfg.GeminiAPIKey)

	notes, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("could not scan notes: %w", err)
	}

	var names []string
	for _, note := range notes {
		names = append(names, note.Title)
	}

	items := []list.Item{}
	for _, note := range notes {
		items = append(items, item{
			title:       note.Title,
			description: note.RelVaultPath,
			Note:        note,
		})
	}

	keys := newListKeyMap()
	list := newNoteList(items, keys)
	m := model{
		list: list,
		keys: keys,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	resultModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	selectedModel, ok := resultModel.(model)
	if ok {
		if selectedModel.selectedItem == nil {
			return fmt.Errorf("no item selected")
		}

		fmt.Println("processing:", selectedModel.selectedItem.title)
		// TODO: add spinner

		if err := transformer.Process(selectedModel.selectedItem.Note); err != nil {
			return fmt.Errorf("could not process note: %w", err)
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
