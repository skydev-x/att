package model

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
