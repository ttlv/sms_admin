package config

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/qor/mailer"
	"github.com/qor/redirect_back"
	"github.com/qor/session/manager"
	"github.com/unrolled/render"
)

var (
	Root         = os.Getenv("GOPATH") + "/src/github.com/ttlv/sms_admin"
	Mailer       *mailer.Mailer
	Render       = render.New()
	_config      *Config
	RedirectBack = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})
)

type Config struct {
	HTTPS            bool
	ServerPort       uint
	DBName           string
	Host             string
	DBPort           string
	User             string
	Password         string
	HttpAuthName     string
	HttpAuthPassword string
}

func MustGetConfig() Config {
	if _config != nil {
		return *_config
	}

	_config = &Config{}
	err := configor.New(&configor.Config{ENVPrefix: "SMSADMIN"}).Load(_config)
	if err != nil {
		panic(err)
	}

	return *_config
}
