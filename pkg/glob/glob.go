package glob

type Config struct {
	ServeAddr string `yaml:"serveAddr"`
}

var G_CONF *Config
