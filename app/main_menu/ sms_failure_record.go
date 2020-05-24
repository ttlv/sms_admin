package main_menu

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/ttlv/sms_admin/config/application"
	"github.com/ttlv/sms_admin/models"
	"github.com/ttlv/sms_admin/utils"
)

func ConfigSmsFailureRecordRes(application *application.Application) {
	res := application.Admin.AddResource(&models.SmsFailureRecord{})
	res.UseTheme("readonly")
	res.SearchAttrs("")
	searchHandler := res.SearchHandler
	res.SearchHandler = func(keyword string, context *qor.Context) *gorm.DB {
		context.SetDB(context.DB.Preload("SmsRecord"))
		return searchHandler(keyword, context)
	}

	res.Filter(&admin.Filter{
		Name:  "Phone",
		Label: "手机号码",
		Handler: func(db *gorm.DB, arg *admin.FilterArgument) *gorm.DB {
			var (
				phoneNumber = arg.Value.Get("Value").Value.([]string)[0]
			)
			_, raw, _ := utils.ParsePhoneNumber(phoneNumber)
			return db.Where("phone = ?", raw)
		},
	})

	res.Meta(&admin.Meta{
		Name:  "CreatedAt",
		Label: "创建时间",
	})

	res.Meta(&admin.Meta{
		Name:  "Phone",
		Label: "手机号",
	})

	res.Meta(&admin.Meta{
		Name:  "ProviderName",
		Label: "短信服务商",
	})

	res.Meta(&admin.Meta{
		Name:  "Error",
		Label: "错误",
	})

	res.IndexAttrs("CreatedAt", "Brand", "Phone", "ProviderName", "Error")
	res.Permission = roles.Deny(roles.Create, roles.Anyone).Deny(roles.Update, roles.Anyone).Deny(roles.Delete, roles.Anyone)
}
