package admin


import (
	"github.com/qor/action_bar"
	"github.com/qor/admin"
	"github.com/qor/media/asset_manager"
	"github.com/qor/media/media_library"
	"github.com/ttlv/sms_admin/config/application"
	"github.com/ttlv/sms_admin/config/i18n"
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

	// Add Media Library
	Admin.AddResource(&media_library.MediaLibrary{}, &admin.Config{Menu: []string{"Site Management"}})

	// Add action bar
	ActionBar = action_bar.New(Admin)
	ActionBar.RegisterAction(&action_bar.Action{Name: "Admin Dashboard", Link: "/admin"})

	// Add Translations
	Admin.AddResource(i18n.I18n, &admin.Config{Menu: []string{"Site Management"}, Priority: -1})

	SetupDashboard(Admin)

	application.Router.Mount(app.Config.Prefix, Admin.NewServeMux(app.Config.Prefix))
}
