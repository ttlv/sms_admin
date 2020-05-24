package main_menu

import (
	"github.com/qor/admin"
	"github.com/ttlv/sms_admin/config/application"
	"github.com/ttlv/sms_admin/models"
)

func ConfigSmsSettingRes(application *application.Application) {
	res := application.Admin.AddResource(&models.SmsSetting{}, &admin.Config{Singleton: true})
	res.EditAttrs("ProviderSorts")
}
