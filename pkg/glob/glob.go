package glob

type Config struct {
	Test      int    `yaml:"test"`
	ServeAddr string `yaml:"serveAddr"`
	DbAddr    string `yaml:"dbAddr"`
	Admin     struct {
		Id string `yaml:"id"`
		Pw string `yaml:"pw"`
	} `yaml:"admin"`
}

var G_CONF *Config
