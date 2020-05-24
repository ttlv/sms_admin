package auth

import (
	"time"

	"github.com/qor/auth"
	"github.com/qor/auth/authority"
	"github.com/qor/auth_themes/clean"
	"github.com/qor/render"
	"github.com/ttlv/sms_admin/config"
	"github.com/ttlv/sms_admin/config/bindatafs"
	"github.com/ttlv/sms_admin/config/db"
	"github.com/ttlv/sms_admin/models/users"
)

var (
	// Auth initialize Auth for Authentication
	Auth = clean.New(&auth.Config{
		DB:         db.DB,
		Mailer:     config.Mailer,
		Render:     render.New(&render.Config{AssetFileSystem: bindatafs.AssetFS.NameSpace("auth")}),
		UserModel:  users.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
	})

	// Authority initialize Authority for Authorization
	Authority = authority.New(&authority.Config{
		Auth: Auth,
	})
)

func init() {
	Authority.Register("logged_in_half_hour", authority.Rule{TimeoutSinceLastLogin: time.Minute * 30})
}
