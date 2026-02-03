# at 🎯

Accidentally Doing Hundred Different-things (ADHD) Tracker; aka `at`

> A simple, Git-backed progress tracker for your daily goals. Track habits, monitor streaks, and stay accountable—all from your terminal.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## ✨ Features

- 📊 **Visual Dashboard** - See your daily progress at a glance
- 🔥 **Streak Tracking** - Monitor consecutive days of activity
- 📝 **Simple Check-ins** - Quick logging with `at checkin`
- 🎯 **Custom Goals** - Set daily targets for each activity
- 🔄 **Git-Backed** - Automatic version control of your progress
- 🚀 **Fast & Lightweight** - Built with Go, runs anywhere
- 🎨 **Emoji Support** - Personalize your topics with emojis

## 🚀 Installation

### Quick Install (macOS & Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/skydev-x/at/master/install.sh | bash
```

### Manual Installation

**Prerequisites:** Go 1.20+

```bash
git clone https://github.com/skydev-x/at.git
cd at
go build -o at .
sudo mv at /usr/local/bin/
```

### Verify Installation

```bash
at --version
```

## 🎯 Quick Start

### 1. Run Setup Wizard

```bash
at setup
```

This will guide you through:

- Creating your data directory
- Setting up Git tracking (optional)
- Configuring your first topics

### 2. Add Your Topics

```bash
# Add topics with: ID, goal description, daily target, and emoji
at topic add dsa "DSA Practice" 3 "💻"
at topic add reading "Daily Reading" 1 "📚"
at topic add workout "Exercise" 1 "💪"
at topic add learning "Learn Something New" 2 "🧠"
```

### 3. Start Tracking

```bash
# Log your activities
at checkin dsa "Solved binary search problems"
at checkin reading "Read chapter 5 of Atomic Habits"
at c workout "30 min cardio"  # 'c' is shorthand for checkin
```

### 4. View Your Dashboard

```bash
at
```

Example output:

```
┌─────────────────────────────────────────────────┐
│           AT - Activity Tracker                 │
│              March 15, 2024                     │
└─────────────────────────────────────────────────┘

💻 DSA Practice          [██░░░] 2/3  🔥 5 days
📚 Daily Reading         [███░░] 1/1  🔥 12 days
💪 Exercise              [░░░░░] 0/1
🧠 Learn Something New   [█░░░░] 1/2  🔥 3 days
```

## 📖 Usage

### Core Commands

| Command                       | Description                          |
| ----------------------------- | ------------------------------------ |
| `at`                          | Show dashboard with today's progress |
| `at checkin <topic> <remark>` | Log an activity                      |
| `at c <topic> <remark>`       | Shorthand for checkin                |
| `at help`                     | Show detailed help                   |

### Topic Management

```bash
# List all topics
at topic list

# Add a new topic (ID, description, daily goal, emoji)
at topic add coding "Coding Practice" 2 "⌨️"

# Pause tracking a topic (without deleting history)
at topic disable coding

# Resume tracking
at topic enable coding

# Remove a topic permanently
at topic remove coding
```

### Configuration

```bash
# View current configuration
at config show

# Change data directory
at config set-path ~/my-tracking-data

# Set Git remote for syncing across devices
at config set-remote git@github.com:yourusername/at-data.git
```

## 💡 Examples

### Daily Routine Tracking

```bash
# Morning routine
at c meditation "10 min mindfulness" 🧘
at c reading "Read for 20 minutes" 📖

# Work tracking
at c coding "Built user authentication" 💻
at c learning "Watched Go tutorial" 🎓

# Evening routine
at c workout "Gym session" 💪
at c journal "Reflected on the day" 📝
```

### Project-Based Tracking

```bash
# Add project-specific topics
at topic add frontend "Frontend Development" 3 "🎨"
at topic add backend "Backend Development" 3 "⚙️"
at topic add testing "Write Tests" 2 "🧪"

# Log your work
at c frontend "Implemented responsive navbar"
at c backend "Created REST API endpoints"
at c testing "Added unit tests for auth module"
```

## 🔧 Configuration Files

### Config Location

`~/.at_config.json`

```json
{
  "data_path": "/Users/you/.at",
  "git_enabled": true,
  "git_remote": "git@github.com:yourusername/at-data.git"
}
```

### Data Directory Structure

```
~/.at/
├── topics.json          # Topic definitions
├── checkins.json        # All check-in records
└── .git/               # Git repository (if enabled)
```

## 🎨 Emoji Ideas

| Category     | Emojis            |
| ------------ | ----------------- |
| **Coding**   | 💻 ⌨️ 🖥️ 👨‍💻 🚀 ⚡ |
| **Learning** | 📚 🎓 🧠 📖 ✍️ 🎯 |
| **Health**   | 💪 🏃 🧘 🥗 💤 🏋️ |
| **Creative** | 🎨 🎵 📸 ✏️ 🎬 🎭 |
| **Personal** | 🌱 💭 ☕ 🎮 🌟 ✨ |

## 🔄 Git Sync (Optional)

Keep your progress synced across multiple devices:

```bash
# 1. Create a private GitHub repo for your data
# 2. Set the remote
at config set-remote git@github.com:yourusername/at-data.git

# Your check-ins will auto-commit and push!
```

## 🛠️ Development

### Build from Source

```bash
git clone https://github.com/skydev-x/at.git
cd at
go build -o at .
```

### Run Tests

```bash
go test ./...
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📊 Use Cases

- **Daily Habit Tracking** - Build consistent routines
- **Learning Goals** - Track study sessions and courses
- **Project Management** - Monitor daily work on projects
- **Health & Fitness** - Log workouts and healthy habits
- **Creative Work** - Track writing, art, or music practice
- **Skill Development** - Monitor practice time for new skills

## ❓ FAQ

**Q: Does this require an internet connection?**  
A: No! Works completely offline. Git sync is optional.

**Q: Can I use this on Windows?**  
A: Currently optimized for macOS and Linux. Windows support coming soon.

**Q: How do I backup my data?**  
A: Your data is in `~/.at/`. Either enable Git sync or manually backup this directory.

**Q: Can I edit past check-ins?**  
A: Currently, you can manually edit `~/.at/checkins.json`. Built-in editing coming soon.

**Q: Is my data private?**  
A: Yes! Everything is stored locally. Git sync is optional and you control the repository.

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with:

- [Go](https://golang.org/) - The programming language
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework (if you're using it)
- Love for productivity and self-improvement ❤️

## 📧 Support

- 🐛 **Bug Reports:** [GitHub Issues](https://github.com/skydev-x/at/issues)
- 💡 **Feature Requests:** [GitHub Discussions](https://github.com/skydev-x/at/discussions)
- 📖 **Documentation:** [Wiki](https://github.com/skydev-x/at/wiki)

---

**Made with ❤️ by [@skydev-x](https://github.com/skydev-x)**

⭐ If you find this helpful, please star the repo!
