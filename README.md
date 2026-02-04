# att 🎯

Accidentally Doing Hundred Different-things (ADHD) Tracker Tool; aka `att`

> A simple, Git-backed progress tracker for your daily goals. Track habits, monitor streaks, and stay accountable—all from your terminal.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## ✨ Features

- 📊 **Visual Dashboard** - See your daily progress at a glance
- 🔥 **Streak Tracking** - Monitor consecutive days of activity
- 📝 **Simple Check-ins** - Quick logging with `att checkin`
- 🎯 **Custom Goals** - Set daily targets for each activity
- 🔄 **Git-Backed** - Automatic version control of your progress
- 🚀 **Fast & Lightweight** - Built with Go, runs anywhere
- 🎨 **Emoji Support** - Personalize your topics with emojis

## 🚀 Installation

### Quick Install (macOS & Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/skydev-x/att/master/install.sh | bash
```

### Manual Installation

**Prerequisites:** Go 1.20+

```bash
git clone https://github.com/skydev-x/att.git
cd att
go build -o att .
sudo mv att /usr/local/bin/
```

### Verify Installation

```bash
att --version
```

## 🎯 Quick Start

### 1. Run Setup Wizard

```bash
att setup
```

This will guide you through:

- Creating your data directory
- Setting up Git tracking (optional)
- Configuring your first topics

### 2. Add Your Topics

```bash
# Add topics with: ID, goal description, daily target, and emoji
att topic add dsa "DSA Practice" 3 "💻"
att topic add reading "Daily Reading" 1 "📚"
att topic add workout "Exercise" 1 "💪"
att topic add learning "Learn Something New" 2 "🧠"
```

### 3. Start Tracking

```bash
# Log your activities
att checkin dsa "Solved binary search problems"
att checkin reading "Read chapter 5 of Atomic Habits"
att c workout "30 min cardio"  # 'c' is shorthand for checkin
```

### 4. View Your Dashboard

```bash
att
```

Example output:

<img width="1626" height="1340" alt="Image" src="https://github.com/user-attachments/assets/372e68b4-3ff7-4a3d-a035-d44ea19ca51b" />
<img width="918" height="520" alt="Image" src="https://github.com/user-attachments/assets/e700df90-8a6c-4ce9-8fce-82d1eba507e8" />

```
💻 DSA Practice          [██░░░] 2/3  🔥 5 days
📚 Daily Reading         [███░░] 1/1  🔥 12 days
💪 Exercise              [░░░░░] 0/1
🧠 Learn Something New   [█░░░░] 1/2  🔥 3 days
```

## 📖 Usage

### Core Commands

| Command                        | Description                          |
| ------------------------------ | ------------------------------------ |
| `att`                          | Show dashboard with today's progress |
| `att checkin <topic> <remark>` | Log an activity                      |
| `att c <topic> <remark>`       | Shorthand for checkin                |
| `att help`                     | Show detailed help                   |

### Topic Management

```bash
# List all topics
att topic list

# Add a new topic (ID, description, daily goal, emoji)
att topic add coding "Coding Practice" 2 "⌨️"

# Pause tracking a topic (without deleting history)
att topic disable coding

# Resume tracking
att topic enable coding

# Remove a topic permanently
att topic remove coding
```

### Configuration

```bash
# View current configuration
att config show

# Change data directory
att config set-path ~/my-tracking-data

# Set Git remote for syncing across devices
att config set-remote git@github.com:yourusername/at-data.git
```

## 💡 Examples

### Daily Routine Tracking

```bash
# Morning routine
att c meditation "10 min mindfulness" 🧘
att c reading "Read for 20 minutes" 📖

# Work tracking
att c coding "Built user authentication" 💻
att c learning "Watched Go tutorial" 🎓

# Evening routine
att c workout "Gym session" 💪
att c journal "Reflected on the day" 📝
```

### Project-Based Tracking

```bash
# Add project-specific topics
att topic add frontend "Frontend Development" 3 "🎨"
att topic add backend "Backend Development" 3 "⚙️"
att topic add testing "Write Tests" 2 "🧪"

# Log your work
att c frontend "Implemented responsive navbar"
att c backend "Created REST API endpoints"
att c testing "Added unit tests for auth module"
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
att config set-remote git@github.com:yourusername/att-data.git

# Your check-ins will auto-commit and push!
```

## 🛠️ Development

### Build from Source

```bash
git clone https://github.com/skydev-x/att.git
cd att
go build -o att .
```

### Run Tests

```bash
go test ./...
```

### TODOs

- [ ] Improve config and git sync
- [ ] Improve viewing of trackings
- [ ] Improve data storage from JSON file to an efficient alternative
- [ ] New commands to view graphs
- [ ] Structure app for better development
- [ ] Add custom styling and layouts

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
A: Your data is in `~/.att/`. Either enable Git sync or manually backup this directory.

**Q: Can I edit past check-ins?**  
A: Currently, you can manually edit `~/.att/checkins.json`. Built-in editing coming soon.

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

- 🐛 **Bug Reports:** [GitHub Issues](https://github.com/skydev-x/att/issues)
- 💡 **Feature Requests:** [GitHub Discussions](https://github.com/skydev-x/att/discussions)
- 📖 **Documentation:** [Wiki](https://github.com/skydev-x/att/wiki)

---

**Made with ❤️ by [@skydev-x](https://github.com/skydev-x)**

⭐ If you find this helpful, please star the repo! or support open source work

<iframe src="https://github.com/sponsors/skydev-x/button" title="Sponsor skydev-x" height="32" width="114" style="border: 0; border-radius: 6px;"></iframe>
