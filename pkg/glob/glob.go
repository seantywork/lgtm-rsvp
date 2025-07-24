package glob

type Config struct {
	ServeAddr    string `yaml:"serveAddr"`
	Url          string `yaml:"url"`
	SessionStore string `yaml:"sessionStore"`
	Db           struct {
		Addr     string `yaml:"addr"`
		InitFile string `yaml:"initFile"`
	} `yaml:"db"`
	Admin struct {
		Id    string `yaml:"id"`
		Pw    string `yaml:"pw"`
		OAuth *struct {
			ClientId                string   `yaml:"clientId"`
			ProjectId               string   `yaml:"projectId"`
			AuthUri                 string   `yaml:"authUri"`
			TokenUri                string   `yaml:"tokenUri"`
			AuthProviderX509CertUrl string   `yaml:"authProviderX509CertUrl"`
			ClientSecret            string   `yaml:"clientSecret"`
			RidirectUris            []string `yaml:"redirectUris"`
		} `yaml:"oauth"`
	} `yaml:"admin"`
	Title    string `yaml:"title"`
	Groom    string `yaml:"groom"`
	Bride    string `yaml:"bride"`
	Comment  string `yaml:"comment"`
	Message  string `yaml:"message"`
	GiftPage string `yaml:"giftPage"`
	Api      struct {
		GoogleComment *string `yaml:"googleComment"`
		KakaoShare    *string `yaml:"kakaoShare"`
	} `yaml:"api"`
}

var G_CONF *Config

var G_CONFIG_PATH = "config/config.yaml"

var G_MEDIA_PATH = "data/media/"

var G_MAIL_ERR_PATH = "data/comment-list.json"

var G_ALBUM_PATH = "public/images/album"

var G_COMMENT_LIST_PATH = "comment-list.json"
