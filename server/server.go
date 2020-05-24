package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/qor/admin"
	"github.com/qor/qor/utils"
	"github.com/ttlv/common_utils/readonly"
	adminapp "github.com/ttlv/sms_admin/app/admin"
	"github.com/ttlv/sms_admin/app/main_menu"
	"github.com/ttlv/sms_admin/app/static"
	"github.com/ttlv/sms_admin/config"
	"github.com/ttlv/sms_admin/config/application"
	"github.com/ttlv/sms_admin/config/bindatafs"
	"github.com/ttlv/sms_admin/config/db"
	"github.com/ttlv/sms_admin/models"
	"net/http"
	"path/filepath"
)

func NewServer() (http.Handler, *admin.Admin) {
	var (
		Router = chi.NewRouter()
		Admin  = admin.New(&admin.AdminConfig{
			SiteName: "SMS ADMIN",
			DB:       db.DB,
		})
		Application = application.New(&application.Config{
			Router: Router,
			Admin:  Admin,
			DB:     db.DB,
		})
	)
	db.DB.AutoMigrate(&models.SmsBrand{}, &models.SmsRecord{}, &models.SmsAvailable{}, &models.SmsFailureRecord{}, &models.SmsSetting{})
	Router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.Header.Del("Authorization")
			handler.ServeHTTP(w, req)
		})
	})
	Router.Use(middleware.RealIP)
	Router.Use(middleware.Logger)
	Router.Use(middleware.Recoverer)
	Application.Use(adminapp.New(&adminapp.Config{}))
	Application.Use(main_menu.New(&main_menu.Config{}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"/system"},
		Handler: utils.FileServer(http.Dir(filepath.Join(config.Root, "public"))),
	}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"javascripts", "stylesheets", "images", "dist", "fonts", "vendors", "favicon.ico"},
		Handler: bindatafs.AssetFS.FileServer(http.Dir("public"), "javascripts", "stylesheets", "images", "dist", "fonts", "vendors", "favicon.ico"),
	}))
	readonly.Setup(Admin)
	return Application.NewServeMux(), Admin
}
