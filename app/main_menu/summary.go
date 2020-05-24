package main_menu

import (
	"github.com/ttlv/sms_admin/config/application"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
	Prefix string
}

// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	ConfigBrandRes(application)
	ConfigSmsAvailableRes(application)
	ConfigSmsFailureRecordRes(application)
	ConfigSmsRecoardRes(application)
	ConfigSmsSettingRes(application)
}
