package config

// Reddit contains configs for reddit
type Reddit struct {
	ClientID   string     `yaml:"clientID"`
	Secret     string     `yaml:"secret"`
	Username   string     `yaml:"username"`
	Password   string     `yaml:"password"`
	Limit      int        `yaml:"limit"`
	Categories []Category `yaml:"categories"`
}

// Category is a categories
type Category struct {
	Name   string   `yaml:"name"`
	Source []string `yaml:"source"`
}
