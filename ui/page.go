package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// Page represents a full-featured terminal page component
type Page struct {
	breadcrumb string
	title      string
	subtitle   string
	content    string
	width      int
	height     int
}

// Styles
var (
	// Color scheme
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	secondaryColor = lipgloss.Color("#6B7280") // Gray
	accentColor    = lipgloss.Color("#10B981") // Green
	textColor      = lipgloss.Color("#F3F4F6") // Light gray
	borderColor    = lipgloss.Color("#374151") // Dark gray

	// Breadcrumb style
	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// Subtitle style
	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Padding(0, 1)

	// Content style
	contentStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(1, 2)

	// Border style
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2)

	// Header separator
	separatorStyle = lipgloss.NewStyle().
			Foreground(borderColor)

	// Footer style
	footerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			Align(lipgloss.Center)
)

// NewPage creates a new page component
func NewPage(breadcrumb, title, subtitle, content string) Page {
	return Page{
		breadcrumb: breadcrumb,
		title:      title,
		subtitle:   subtitle,
		content:    content,
		width:      80,
		height:     24,
	}
}

func (p Page) Init() tea.Cmd {
	return nil
}

func (p Page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p Page) View() string {
	// Calculate dimensions
	contentWidth := p.width - 6 // Account for border and padding

	// Build breadcrumb
	breadcrumb := ""
	if p.breadcrumb != "" {
		breadcrumb = breadcrumbStyle.Render("  " + p.breadcrumb)
	}

	// Build title
	title := titleStyle.Width(contentWidth).Render(p.title)

	// Build subtitle
	subtitle := ""
	if p.subtitle != "" {
		subtitle = subtitleStyle.Width(contentWidth).Render(p.subtitle)
	}

	// Build separator
	separator := separatorStyle.Render(strings.Repeat("─", contentWidth))

	// Build content
	content := contentStyle.
		Width(contentWidth).
		Render(p.content)

	// Build footer
	footer := footerStyle.
		Width(contentWidth).
		Render("Press 'q' or Ctrl+C to quit")

	// Combine all parts
	var parts []string

	if breadcrumb != "" {
		parts = append(parts, breadcrumb)
	}
	parts = append(parts, title)
	if subtitle != "" {
		parts = append(parts, subtitle)
	}
	parts = append(parts, separator)
	parts = append(parts, content)
	parts = append(parts, "")
	parts = append(parts, footer)

	page := lipgloss.JoinVertical(lipgloss.Left, parts...)

	// Add border
	bordered := borderStyle.
		Width(contentWidth).
		Render(page)

	return bordered
}
