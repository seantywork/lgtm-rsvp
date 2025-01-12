package glob

type Config struct {
	Test         int    `yaml:"test"`
	ServeAddr    string `yaml:"serveAddr"`
	Url          string `yaml:"url"`
	SessionStore string `yaml:"sessionStore"`
	Db           struct {
		Addr     string `yaml:"addr"`
		InitFile string `yaml:"initFile"`
	} `yaml:"db"`
	Admin struct {
		UseOauth2 bool   `yaml:"useOauth2"`
		Id        string `yaml:"id"`
		Pw        string `yaml:"pw"`
	} `yaml:"admin"`
}

var G_CONF *Config
