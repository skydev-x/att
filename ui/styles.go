package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			PaddingLeft(1)

	TopicStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true).
			MarginTop(1)

	DisabledStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Strikethrough(true)

	StatsStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			PaddingLeft(2)

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2)

	FooterStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			Align(lipgloss.Center).
			MarginTop(1)

	SuccessBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Padding(1, 2)

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DangerColor).
			Padding(1, 2)
)
