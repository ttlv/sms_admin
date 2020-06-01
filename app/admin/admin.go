package admin

import (
	"github.com/qor/action_bar"
	"github.com/qor/admin"
	"github.com/qor/media/asset_manager"
	"github.com/ttlv/sms_admin/config/application"
)

// ActionBar admin action bar
var ActionBar *action_bar.ActionBar

// AssetManager asset manager
var AssetManager *admin.Resource

// New new home app
func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/admin"
	}
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
	Admin := application.Admin
	AssetManager = Admin.AddResource(&asset_manager.AssetManager{}, &admin.Config{Invisible: true})
	application.Router.Mount(app.Config.Prefix, Admin.NewServeMux(app.Config.Prefix))
}
