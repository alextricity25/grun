package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

/*
This example assumes an existing understanding of commands and messages. If you
haven't already read our tutorials on the basics of Bubble Tea and working
with commands, we recommend reading those first.

Find them at:
https://github.com/charmbracelet/bubbletea/tree/master/tutorials/commands
https://github.com/charmbracelet/bubbletea/tree/master/tutorials/basics
*/

// sessionState is used to track which model is focused
type sessionState uint

type item string

type itemDelegate struct{}

func (i item) FilterValue() string { return "" }

const (
	defaultTime              = time.Minute
	timerView   sessionState = iota
	spinnerView
	tabView
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}
	modelStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))
	focusedModelStyle = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	spinnerStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 1, 1, 1)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type mainModel struct {
	state              sessionState
	timer              timer.Model
	spinner            spinner.Model
	index              int
	Tabs               []string
	TabContent         []string
	PreviewPaneContent string
	activeTab          int
	grunServices       list.Model
}
func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func newModel(timeout time.Duration) mainModel {
	m := mainModel{state: tabView}
	m.timer = timer.New(timeout)
	m.spinner = spinner.New()
	m.Tabs = []string{
		"Services",
		"Jobs",
		"Info",
	}

	// grunServices list
	items := []list.Item{}
	for _, service := range getGrunServices() {
		items = append(items, item(service))
	}
	l := list.New(items, itemDelegate{}, 20, 14)
	l.Title = "Cloud Run Services"
	m.grunServices = l
	m.TabContent = []string{
		m.grunServices.View(),
		"Listing Cloud Run Jobs",
		"Listing Info",
	}
	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(m.timer.Init(), m.spinner.Tick)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.grunServices, cmd = m.grunServices.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
		case "left", "h":
			m.activeTab = max(m.activeTab-1, 0)
		case "tab":
			if m.state == timerView {
				m.state = spinnerView
			} else {
				m.state = timerView
			}
		case "n":
			if m.state == timerView {
				m.timer = timer.New(defaultTime)
				cmds = append(cmds, m.timer.Init())
			} else {
				m.Next()
				m.resetSpinner()
				cmds = append(cmds, m.spinner.Tick)
			}
		}
		switch m.state {
		// update whichever model is focused
		case spinnerView:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.timer, cmd = m.timer.Update(msg)
			cmds = append(cmds, cmd)
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) TabView() string {
	doc := strings.Builder{}
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	tabTableWidth := (physicalWidth / 2)
	realTabTableWidth := ((tabTableWidth / len(m.Tabs)) * len(m.Tabs))

	var renderedTabs []string
	var tabContent string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		style.Width((realTabTableWidth / len(m.Tabs)) - 1)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	if m.activeTab == 0 {
		tabContent = m.grunServices.View()
	} else {
		tabContent = m.TabContent[m.activeTab]
	}
	doc.WriteString(windowStyle.Width((realTabTableWidth + 1)).Render(tabContent))
	return docStyle.Render(doc.String())
}

func (m mainModel) PreviewPaneView() string {
	doc := strings.Builder{}
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	previewPaneView := (physicalWidth / 2) - 15
	style := activeTabStyle.Copy()

	style.Width(previewPaneView).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))

	if m.activeTab == 0 {
		doc.WriteString(style.Render("Listing Cloud Run Services"))
	} else if m.activeTab == 1 {
		doc.WriteString(style.Render("Listing Cloud Run Jobs"))
	} else {
		doc.WriteString(style.Render("Listing Info What if this line wraps is super long like it's so long that it needs to wrap all the way"))
	}
	return docStyle.Render(doc.String())
}

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()
	if m.state == tabView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, m.TabView(), m.PreviewPaneView())
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), focusedModelStyle.Render(m.spinner.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == timerView {
		return "timer"
	}
	return "spinner"
}

func (m *mainModel) Next() {
	if m.index == len(spinners)-1 {
		m.index = 0
	} else {
		m.index++
	}
}

func (m *mainModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}
