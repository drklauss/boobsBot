package config

// Telegram contains configs for telegram.
type Telegram struct {
	Workers int `yaml:"workers"`
	Proxy   *struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user,omitempty"`
		Password string `yaml:"password,omitempty"`
	} `yaml:"proxy,omitempty"`
	Token   string `yaml:"token"`
	API     string `yaml:"api"`
	BotName string `yaml:"bot_name"`
	Admin   []int  `yaml:"admin"`
	Time    struct {
		Update            int   `yaml:"update"`
		SkipMessages      int64 `yaml:"skip_messages"`
		QuerySkipMessages int64 `yaml:"query_skip_messages"`
	} `yaml:"time,flow"`
}
