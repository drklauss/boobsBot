package config

// Telegram contains configs for telegram
type Telegram struct {
	Workers int `yaml:"workers"`
	Proxy   *struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user,omitempty"`
		Password string `yaml:"password,omitempty"`
	} `yaml:"proxy,omitempty"`
	Token string `yaml:"token"`
	API   string `yaml:"api"`
	Admin []int  `yaml:"admin"`
	Time  struct {
		Update       int `yaml:"update"`
		SkipMessages int `yaml:"skip_messages"`
	} `yaml:"time,flow"`
}
