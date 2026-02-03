package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	primaryColor = lipgloss.Color("#7C3AED")
	successColor = lipgloss.Color("#10B981")
	warningColor = lipgloss.Color("#F59E0B")
	dangerColor  = lipgloss.Color("#EF4444")
	mutedColor   = lipgloss.Color("#6B7280")
	textColor    = lipgloss.Color("#F3F4F6")
	borderColor  = lipgloss.Color("#374151")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			PaddingLeft(1)

	topicStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true).
			MarginTop(1)

	disabledStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Strikethrough(true)

	statsStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			PaddingLeft(2)

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Align(lipgloss.Center).
			MarginTop(1)

	successBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(successColor).
			Padding(1, 2)

	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dangerColor).
			Padding(1, 2)
)

// Config structures
type TopicConfig struct {
	Name      string `json:"name"`
	DailyGoal int    `json:"daily_goal"`
	Emoji     string `json:"emoji"`
	Enabled   bool   `json:"enabled"`
}

type Config struct {
	DataPath string                  `json:"data_path"`
	SSHURL   string                  `json:"ssh_url,omitempty"`
	Topics   map[string]*TopicConfig `json:"topics"`
}

// Data structures
type CheckIn struct {
	Date   string `json:"date"`
	Remark string `json:"remark"`
}

type TopicData struct {
	Name          string    `json:"name"`
	TotalCheckIns int       `json:"total_checkins"`
	Streak        int       `json:"streak"`
	LastDate      string    `json:"last_date,omitempty"`
	History       []CheckIn `json:"history"`
}

type ProgressData struct {
	Created string                `json:"created"`
	Topics  map[string]*TopicData `json:"topics"`
}

// Dashboard model
type Dashboard struct {
	cfg    *Config
	data   *ProgressData
	width  int
	height int
	err    error
}

func main() {

	VERSION := "v1.0.0"
	if len(os.Args) == 1 {
		showDashboard()
		return
	}

	command := os.Args[1]

	switch command {
	case "checkin", "c":
		if len(os.Args) < 4 {
			fmt.Println("Usage: att checkin <topic> <remark>")
			os.Exit(1)
		}
		checkin(os.Args[2], strings.Join(os.Args[3:], " "))
	case "topic":
		handleTopicCommand()
	case "config":
		handleConfigCommand()
	case "setup":
		runSetup()
	case "help", "--help", "-h":
		showHelp()
	case "version", "--version", "-v":
		fmt.Println(VERSION)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'att help' for usage")
		os.Exit(1)
	}
}

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".att_config.json")
}

func getDefaultDataPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".att")
}

func loadConfig() *Config {
	configPath := getConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No config exists - need initial setup
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	// Initialize topics map if nil
	if cfg.Topics == nil {
		cfg.Topics = make(map[string]*TopicConfig)
	}

	return &cfg
}

func saveConfig(cfg *Config) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(getConfigPath(), data, 0644); err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		os.Exit(1)
	}
}

func runGit(repoPath string, args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", repoPath}, args...)...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func initRepo(cfg *Config) {
	dataPath := cfg.DataPath
	os.MkdirAll(dataPath, 0755)

	gitDir := filepath.Join(dataPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return // Already initialized
	}

	runGit(dataPath, "init")

	// Create initial data file
	data := &ProgressData{
		Created: time.Now().Format(time.RFC3339),
		Topics:  make(map[string]*TopicData),
	}

	for topicID, topicCfg := range cfg.Topics {
		data.Topics[topicID] = &TopicData{
			Name:    topicCfg.Name,
			History: []CheckIn{},
		}
	}

	saveData(dataPath, data)
	runGit(dataPath, "add", ".")
	runGit(dataPath, "commit", "-m", "Initial commit")

	if cfg.SSHURL != "" {
		runGit(dataPath, "remote", "add", "origin", cfg.SSHURL)
	}
}

func syncRepo(dataPath string) {
	runGit(dataPath, "pull", "origin", "main", "--rebase")
	runGit(dataPath, "push", "origin", "main")
}

func loadData(dataPath string) *ProgressData {
	progressPath := filepath.Join(dataPath, "progress.json")

	data, err := os.ReadFile(progressPath)
	if err != nil {
		// Create initial data if doesn't exist
		progressData := &ProgressData{
			Created: time.Now().Format(time.RFC3339),
			Topics:  make(map[string]*TopicData),
		}
		saveData(dataPath, progressData)
		return progressData
	}

	var progressData ProgressData
	if err := json.Unmarshal(data, &progressData); err != nil {
		fmt.Printf("Error parsing data: %v\n", err)
		os.Exit(1)
	}

	return &progressData
}

func saveData(dataPath string, data *ProgressData) {
	progressPath := filepath.Join(dataPath, "progress.json")

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(progressPath, jsonData, 0644); err != nil {
		fmt.Printf("Error writing data: %v\n", err)
		os.Exit(1)
	}

	runGit(dataPath, "add", "progress.json")
	commitMsg := fmt.Sprintf("Update: %s", time.Now().Format("2006-01-02 15:04"))
	runGit(dataPath, "commit", "-m", commitMsg)
}

func checkStreaks(data *ProgressData) {
	today := time.Now().Truncate(24 * time.Hour)

	for _, topicData := range data.Topics {
		if topicData.LastDate != "" {
			lastDate, _ := time.Parse(time.RFC3339, topicData.LastDate)
			lastDate = lastDate.Truncate(24 * time.Hour)
			yesterday := today.AddDate(0, 0, -1)

			if lastDate.Before(yesterday) {
				topicData.Streak = 0
			}
		}
	}
}

func getTodayProgress(data *ProgressData, topicID string) int {
	today := time.Now().Truncate(24 * time.Hour)
	topicData := data.Topics[topicID]
	if topicData == nil {
		return 0
	}

	count := 0
	for _, entry := range topicData.History {
		entryDate, _ := time.Parse(time.RFC3339, entry.Date)
		if entryDate.Truncate(24 * time.Hour).Equal(today) {
			count++
		}
	}
	return count
}

// Dashboard UI
func NewDashboard() *Dashboard {
	cfg := loadConfig()
	if cfg == nil {
		return &Dashboard{err: fmt.Errorf("no configuration found - run 'att setup' first")}
	}

	initRepo(cfg)
	data := loadData(cfg.DataPath)
	checkStreaks(data)

	if cfg.SSHURL != "" {
		syncRepo(cfg.DataPath)
	}

	return &Dashboard{
		cfg:    cfg,
		data:   data,
		width:  80,
		height: 24,
	}
}

func (d *Dashboard) Init() tea.Cmd {
	return nil
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		return d, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return d, tea.Quit
		}
	}
	return d, nil
}

func (d *Dashboard) View() string {
	if d.err != nil {
		return errorBoxStyle.Render(fmt.Sprintf("Error: %v\n\nRun 'att setup' to get started.", d.err))
	}

	contentWidth := d.width - 8
	if contentWidth < 40 {
		contentWidth = 40
	}

	// Title
	title := titleStyle.Render("🚀 Activity Tracker")
	subtitle := subtitleStyle.Render(" Dashboard")

	// Separator
	separator := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat("─", contentWidth))

	var sections []string
	sections = append(sections, title, subtitle, separator)

	// Check if no topics
	if len(d.cfg.Topics) == 0 {
		noTopics := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginTop(1).
			Render("No topics configured. Run 'att topic add' to create your first topic.")
		sections = append(sections, "", noTopics)
	} else {
		// Show topics
		for topicID, topicCfg := range d.cfg.Topics {
			topicData := d.data.Topics[topicID]
			if topicData == nil {
				continue
			}

			// Skip disabled topics
			if !topicCfg.Enabled {
				disabledText := disabledStyle.Render(fmt.Sprintf("%s %s (disabled)", topicCfg.Emoji, topicCfg.Name))
				sections = append(sections, "", disabledText)
				continue
			}

			// Topic header
			topicHeader := topicStyle.Render(fmt.Sprintf("%s %s", topicCfg.Emoji, topicCfg.Name))

			// Progress
			todayProgress := getTodayProgress(d.data, topicID)
			dailyGoal := topicCfg.DailyGoal

			progressBar := ""
			for i := 0; i < dailyGoal; i++ {
				if i < todayProgress {
					progressBar += "█"
				} else {
					progressBar += "░"
				}
			}

			progressText := ""
			if todayProgress >= dailyGoal {
				progressText = lipgloss.NewStyle().
					Foreground(successColor).
					Render(fmt.Sprintf("  Today: %d/%d [%s] ✓ GOAL MET!", todayProgress, dailyGoal, progressBar))
			} else {
				progressText = statsStyle.Render(fmt.Sprintf("  Today: %d/%d [%s]", todayProgress, dailyGoal, progressBar))
			}

			// Streak
			streakText := ""
			if topicData.Streak > 0 {
				streakText = statsStyle.Render("  Streak: ") +
					lipgloss.NewStyle().Foreground(warningColor).Bold(true).
						Render(fmt.Sprintf("%d days 🔥", topicData.Streak))
			} else {
				streakText = statsStyle.Render("  Streak: 0 days (start today!)")
			}

			// Total
			totalText := statsStyle.Render(fmt.Sprintf("  Total: %d check-ins", topicData.TotalCheckIns))

			// Today's check-ins
			today := time.Now().Truncate(24 * time.Hour)
			var todayEntries []CheckIn
			for _, entry := range topicData.History {
				entryDate, _ := time.Parse(time.RFC3339, entry.Date)
				if entryDate.Truncate(24 * time.Hour).Equal(today) {
					todayEntries = append(todayEntries, entry)
				}
			}

			checkInsText := ""
			if len(todayEntries) > 0 {
				checkInsText = statsStyle.Render("  Today's work:") + "\n"
				for _, ci := range todayEntries {
					entryTime, _ := time.Parse(time.RFC3339, ci.Date)
					timeStr := entryTime.Format("15:04")
					checkInsText += lipgloss.NewStyle().
						Foreground(mutedColor).
						PaddingLeft(4).
						Render(fmt.Sprintf("[%s] %s\n", timeStr, ci.Remark))
				}
			}

			sections = append(sections, "", topicHeader, progressText, streakText, totalText)
			if checkInsText != "" {
				sections = append(sections, checkInsText)
			}
		}
	}

	// Footer
	footer := footerStyle.Width(contentWidth).
		Render("Press 'q' to quit • 'att help' for commands")
	sections = append(sections, "", footer)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return borderStyle.Width(contentWidth).Render(content)
}

func showDashboard() {
	p := tea.NewProgram(NewDashboard(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Checkin command
func checkin(topicID, remark string) {
	cfg := loadConfig()
	if cfg == nil {
		fmt.Println("No configuration found. Run 'att setup' first.")
		os.Exit(1)
	}

	topicCfg, exists := cfg.Topics[topicID]
	if !exists {
		fmt.Printf("Unknown topic: %s\n", topicID)
		fmt.Println("\nAvailable topics:")
		for id, tc := range cfg.Topics {
			status := "enabled"
			if !tc.Enabled {
				status = "disabled"
			}
			fmt.Printf("  %s - %s %s (%s)\n", id, tc.Emoji, tc.Name, status)
		}
		os.Exit(1)
	}

	if !topicCfg.Enabled {
		fmt.Printf("Topic '%s' is disabled. Enable it with: att topic enable %s\n", topicID, topicID)
		os.Exit(1)
	}

	initRepo(cfg)
	data := loadData(cfg.DataPath)
	checkStreaks(data)

	if cfg.SSHURL != "" {
		syncRepo(cfg.DataPath)
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	topicData := data.Topics[topicID]
	if topicData == nil {
		topicData = &TopicData{
			Name:    topicCfg.Name,
			History: []CheckIn{},
		}
		data.Topics[topicID] = topicData
	}

	// Add check-in
	topicData.History = append(topicData.History, CheckIn{
		Date:   now.Format(time.RFC3339),
		Remark: remark,
	})
	topicData.TotalCheckIns++

	// Update streak
	currentProgress := getTodayProgress(data, topicID)
	if currentProgress == 1 {
		if topicData.LastDate != "" {
			lastDate, _ := time.Parse(time.RFC3339, topicData.LastDate)
			lastDate = lastDate.Truncate(24 * time.Hour)
			yesterday := today.AddDate(0, 0, -1)

			if lastDate.Equal(yesterday) || lastDate.Equal(today) {
				topicData.Streak++
			} else {
				topicData.Streak = 1
			}
		} else {
			topicData.Streak = 1
		}
	}
	topicData.LastDate = now.Format(time.RFC3339)

	saveData(cfg.DataPath, data)

	if cfg.SSHURL != "" {
		syncRepo(cfg.DataPath)
	}

	// Show success message
	showCheckinSuccess(topicCfg, topicData, currentProgress, remark)
}

func showCheckinSuccess(cfg *TopicConfig, data *TopicData, progress int, remark string) {
	title := lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true).
		Render("✓ Check-in Recorded!")

	topicLine := lipgloss.NewStyle().
		Foreground(textColor).
		Bold(true).
		MarginTop(1).
		Render(fmt.Sprintf("%s %s", cfg.Emoji, cfg.Name))

	remarkLine := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Render(fmt.Sprintf("\"%s\"", remark))

	progressBar := ""
	for i := 0; i < cfg.DailyGoal; i++ {
		if i < progress {
			progressBar += "█"
		} else {
			progressBar += "░"
		}
	}

	progressLine := ""
	if progress >= cfg.DailyGoal {
		progressLine = lipgloss.NewStyle().
			Foreground(successColor).
			Render(fmt.Sprintf("Progress: %d/%d [%s] 🎉", progress, cfg.DailyGoal, progressBar))
	} else {
		progressLine = fmt.Sprintf("Progress: %d/%d [%s]", progress, cfg.DailyGoal, progressBar)
	}

	streakLine := lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true).
		Render(fmt.Sprintf("Streak: %d days 🔥", data.Streak))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		topicLine,
		remarkLine,
		"",
		progressLine,
		streakLine,
	)

	fmt.Println()
	fmt.Println(successBoxStyle.Render(content))
	fmt.Println()
}

// Topic management
func handleTopicCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: att topic <command> [args]")
		fmt.Println("\nCommands:")
		fmt.Println("  add <id> <n> <goal> [emoji]  - Add new topic")
		fmt.Println("  remove <id>                   - Remove topic")
		fmt.Println("  enable <id>                   - Enable topic")
		fmt.Println("  disable <id>                  - Disable topic")
		fmt.Println("  list                          - List all topics")
		os.Exit(1)
	}

	subCmd := os.Args[2]

	switch subCmd {
	case "add":
		topicAdd()
	case "remove", "rm", "delete":
		topicRemove()
	case "enable":
		topicEnable(true)
	case "disable":
		topicDisable(false)
	case "list", "ls":
		topicList()
	default:
		fmt.Printf("Unknown topic command: %s\n", subCmd)
		os.Exit(1)
	}
}

func topicAdd() {
	if len(os.Args) < 6 {
		fmt.Println("Usage: att topic add <id> <n> <goal> [emoji]")
		fmt.Println("\nExamples:")
		fmt.Println("  att topic add dsa 'DSA Practice' 3 '💻'")
		fmt.Println("  att topic add reading 'Daily Reading' 1 '📚'")
		fmt.Println("  att topic add exercise 'Exercise' 1 '💪'")
		os.Exit(1)
	}

	topicID := os.Args[3]
	name := os.Args[4]
	var dailyGoal int
	fmt.Sscanf(os.Args[5], "%d", &dailyGoal)

	emoji := "📌"
	if len(os.Args) > 6 {
		emoji = os.Args[6]
	}

	cfg := loadConfig()
	if cfg == nil {
		cfg = &Config{
			DataPath: getDefaultDataPath(),
			Topics:   make(map[string]*TopicConfig),
		}
	}

	if _, exists := cfg.Topics[topicID]; exists {
		fmt.Printf("Topic '%s' already exists\n", topicID)
		os.Exit(1)
	}

	cfg.Topics[topicID] = &TopicConfig{
		Name:      name,
		DailyGoal: dailyGoal,
		Emoji:     emoji,
		Enabled:   true,
	}

	saveConfig(cfg)

	// Update data if repo exists
	if _, err := os.Stat(filepath.Join(cfg.DataPath, ".git")); err == nil {
		data := loadData(cfg.DataPath)
		data.Topics[topicID] = &TopicData{
			Name:    name,
			History: []CheckIn{},
		}
		saveData(cfg.DataPath, data)

		if cfg.SSHURL != "" {
			syncRepo(cfg.DataPath)
		}
	} else {
		initRepo(cfg)
	}

	fmt.Printf("✓ Topic added: %s %s (goal: %d/day)\n", emoji, name, dailyGoal)
}

func topicRemove() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: att topic remove <id>")
		os.Exit(1)
	}

	topicID := os.Args[3]
	cfg := loadConfig()
	if cfg == nil {
		fmt.Println("No configuration found")
		os.Exit(1)
	}

	if _, exists := cfg.Topics[topicID]; !exists {
		fmt.Printf("Topic '%s' not found\n", topicID)
		os.Exit(1)
	}

	fmt.Printf("Remove topic '%s'? This will delete all history. (y/n): ", topicID)
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "yes" {
		fmt.Println("Cancelled")
		return
	}

	delete(cfg.Topics, topicID)
	saveConfig(cfg)

	// Remove from data
	if _, err := os.Stat(filepath.Join(cfg.DataPath, ".git")); err == nil {
		data := loadData(cfg.DataPath)
		delete(data.Topics, topicID)
		saveData(cfg.DataPath, data)

		if cfg.SSHURL != "" {
			syncRepo(cfg.DataPath)
		}
	}

	fmt.Printf("✓ Topic '%s' removed\n", topicID)
}

func topicEnable(enable bool) {
	if len(os.Args) < 4 {
		action := "enable"
		if !enable {
			action = "disable"
		}
		fmt.Printf("Usage: att topic %s <id>\n", action)
		os.Exit(1)
	}

	topicID := os.Args[3]
	cfg := loadConfig()
	if cfg == nil {
		fmt.Println("No configuration found")
		os.Exit(1)
	}

	topicCfg, exists := cfg.Topics[topicID]
	if !exists {
		fmt.Printf("Topic '%s' not found\n", topicID)
		os.Exit(1)
	}

	topicCfg.Enabled = enable
	saveConfig(cfg)

	status := "enabled"
	if !enable {
		status = "disabled"
	}
	fmt.Printf("✓ Topic '%s' %s\n", topicID, status)
}

func topicDisable(enable bool) {
	topicEnable(false)
}

func topicList() {
	cfg := loadConfig()
	if cfg == nil {
		fmt.Println("No configuration found")
		os.Exit(1)
	}

	if len(cfg.Topics) == 0 {
		fmt.Println("No topics configured")
		fmt.Println("\nAdd a topic with: att topic add <id> <n> <goal> [emoji]")
		return
	}

	fmt.Println("\nConfigured Topics:")
	fmt.Println(strings.Repeat("─", 60))

	for id, topic := range cfg.Topics {
		status := "✓"
		statusColor := successColor
		if !topic.Enabled {
			status = "✗"
			statusColor = dangerColor
		}

		statusText := lipgloss.NewStyle().Foreground(statusColor).Render(status)
		fmt.Printf("%s %s %s - %s (goal: %d/day)\n",
			statusText, id, topic.Emoji, topic.Name, topic.DailyGoal)
	}
	fmt.Println()
}

// Config management
func handleConfigCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: att config <command>")
		fmt.Println("\nCommands:")
		fmt.Println("  show              - Show current configuration")
		fmt.Println("  set-path <path>   - Set data directory path")
		fmt.Println("  set-remote <url>  - Set Git remote URL")
		os.Exit(1)
	}

	subCmd := os.Args[2]

	switch subCmd {
	case "show":
		configShow()
	case "set-path":
		configSetPath()
	case "set-remote":
		configSetRemote()
	default:
		fmt.Printf("Unknown config command: %s\n", subCmd)
		os.Exit(1)
	}
}

func configShow() {
	cfg := loadConfig()
	if cfg == nil {
		fmt.Println("No configuration found. Run 'att setup' to create one.")
		return
	}

	fmt.Println("\n📋 Current Configuration")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("Data Path:   %s\n", cfg.DataPath)

	if cfg.SSHURL != "" {
		fmt.Printf("Git Remote:  %s\n", cfg.SSHURL)
	} else {
		fmt.Println("Git Remote:  Not configured")
	}

	fmt.Printf("Topics:      %d configured\n", len(cfg.Topics))
	fmt.Println()

	if len(cfg.Topics) > 0 {
		fmt.Println("Topics:")
		for id, topic := range cfg.Topics {
			status := "(enabled)"
			if !topic.Enabled {
				status = "(disabled)"
			}
			fmt.Printf("  %s - %s %s %s\n", id, topic.Emoji, topic.Name, status)
		}
		fmt.Println()
	}

	fmt.Printf("Config file: %s\n", getConfigPath())
	fmt.Println()
}

func configSetPath() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: at config set-path <path>")
		fmt.Println("\nExamples:")
		fmt.Println("  att config set-path ~/.att")
		fmt.Println("  att config set-path ~/Documents/tracking")
		os.Exit(1)
	}

	newPath := os.Args[3]

	// Expand ~ to home directory
	if strings.HasPrefix(newPath, "~") {
		home, _ := os.UserHomeDir()
		newPath = filepath.Join(home, newPath[1:])
	}

	cfg := loadConfig()
	if cfg == nil {
		cfg = &Config{
			DataPath: newPath,
			Topics:   make(map[string]*TopicConfig),
		}
	} else {
		cfg.DataPath = newPath
	}

	saveConfig(cfg)
	fmt.Printf("✓ Data path updated: %s\n", newPath)
}

func configSetRemote() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: att config set-remote <url>")
		fmt.Println("\nExamples:")
		fmt.Println("  att config set-remote git@github.com:user/tracking.git")
		fmt.Println("  att config set-remote ''  (to remove)")
		os.Exit(1)
	}

	remoteURL := os.Args[3]

	cfg := loadConfig()
	if cfg == nil {
		fmt.Println("No configuration found. Run 'att setup' first.")
		os.Exit(1)
	}

	cfg.SSHURL = remoteURL
	saveConfig(cfg)

	if remoteURL == "" {
		fmt.Println("✓ Git remote removed")
	} else {
		fmt.Printf("✓ Git remote updated: %s\n", remoteURL)

		// Update git remote if repo exists
		if _, err := os.Stat(filepath.Join(cfg.DataPath, ".git")); err == nil {
			runGit(cfg.DataPath, "remote", "remove", "origin")
			runGit(cfg.DataPath, "remote", "add", "origin", remoteURL)
		}
	}
}

// Setup wizard
func runSetup() {
	fmt.Println()
	fmt.Println("🎯 Activity Tracker Setup")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	cfg := loadConfig()
	if cfg == nil {
		cfg = &Config{
			DataPath: getDefaultDataPath(),
			Topics:   make(map[string]*TopicConfig),
		}
	}

	// Data path
	fmt.Printf("Data directory [%s]: ", cfg.DataPath)
	var dataPath string
	fmt.Scanln(&dataPath)
	if dataPath != "" {
		if strings.HasPrefix(dataPath, "~") {
			home, _ := os.UserHomeDir()
			dataPath = filepath.Join(home, dataPath[1:])
		}
		cfg.DataPath = dataPath
	}

	// Git remote
	currentRemote := cfg.SSHURL
	if currentRemote == "" {
		currentRemote = "none"
	}
	fmt.Printf("Git remote URL [%s]: ", currentRemote)
	var remoteURL string
	fmt.Scanln(&remoteURL)
	if remoteURL != "" {
		cfg.SSHURL = remoteURL
	}

	// Save config
	saveConfig(cfg)

	// Initialize repo
	initRepo(cfg)

	fmt.Println()
	fmt.Println("✓ Setup complete!")
	fmt.Println()

	if len(cfg.Topics) == 0 {
		fmt.Println("No topics configured yet. Add your first topic:")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  att topic add dsa 'DSA Practice' 3 '💻'")
		fmt.Println("  att topic add reading 'Daily Reading' 1 '📚'")
		fmt.Println("  att topic add exercise 'Exercise' 1 '💪'")
		fmt.Println("  att topic add coding 'Coding Projects' 2 '⌨️'")
		fmt.Println()
	} else {
		fmt.Println("Run 'att' to see your dashboard")
	}
}

func showHelp() {
	help := `
AHDHD - Tracker Tool

A Git-backed progress tracker for your daily goals.

USAGE:
  att                                  Show dashboard
  att checkin <topic> <remark>         Record a check-in
  att topic <command> [args]           Manage topics
  att config <command> [args]          Manage configuration
  att setup                            Run setup wizard
  att help                             Show this help

TOPIC COMMANDS:
  att topic add <id> <n> <goal> [emoji]   Add new topic
  att topic remove <id>                    Remove topic
  att topic enable <id>                    Enable topic
  att topic disable <id>                   Disable topic (pause tracking)
  att topic list                           List all topics

CONFIG COMMANDS:
  att config show                      Show configuration
  att config set-path <path>           Set data directory
  att config set-remote <url>          Set Git remote URL

EXAMPLES:
  # Add topics
  att topic add dsa 'DSA Practice' 3 '💻'
  att topic add reading 'Daily Reading' 1 '📚'

  # Check in
  att checkin dsa "Solved two sum problem"
  att c reading "Read 30 pages"  # 'c' is short for checkin

  # Manage topics
  att topic disable dsa              # Pause tracking
  att topic enable dsa               # Resume tracking
  att topic remove old-topic         # Delete topic

  # Configuration
  att config show                    # View settings
  att config set-path ~/tracking     # Change data location
  att config set-remote git@...      # Set Git remote

FILES:
  Config:  ~/.att_config.json
  Data:    ~/.att/ (or custom path)

For more info, visit: (https://github.com/skydev-x/att)`
	fmt.Println(help)
}
