package glob

type Config struct {
	ServeAddr string `yaml:"serveAddr"`
	DbAddr string `yaml:"dbAddr"`
}

var G_CONF *Config
