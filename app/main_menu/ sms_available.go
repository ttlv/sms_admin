package main_menu

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/ttlv/sms_admin/config/application"
	"github.com/ttlv/sms_admin/models"
)

func ConfigSmsAvailableRes(application *application.Application) {
	res := application.Admin.AddResource(&models.SmsAvailable{})
	res.UseTheme("readonly")
	res.SearchAttrs("")
	searchHandler := res.SearchHandler
	res.SearchHandler = func(keyword string, context *qor.Context) *gorm.DB {
		context.SetDB(context.DB.Preload("SmsBrand"))
		return searchHandler(keyword, context)
	}
	res.Meta(&admin.Meta{Name: "BrandName", FieldName: "SmsBrand.Name"})
	res.Meta(&admin.Meta{
		Name: "SmsBranID",
		Config: &admin.SelectOneConfig{
			Collection: func(i interface{}, c *qor.Context) (result [][]string) {
				var (
					brands = []models.SmsBrand{}
				)
				c.DB.Find(&brands)
				for _, brand := range brands {
					result = append(result, []string{fmt.Sprintf("%v", brand.ID), brand.Name})
				}
				return
			},
		},
	})
	res.Meta(&admin.Meta{Name: "Note", Type: "text"})
	res.IndexAttrs("CreatedAt", "BrandName", "AvailableAmount", "Note")
	res.NewAttrs(&admin.Section{
		Rows: [][]string{
			{"SmsBranID"},
			{"AvailableAmount"},
			{"Note"},
		},
	})
	res.Permission = roles.Deny(roles.Update, roles.Anyone).Deny(roles.Delete, roles.Anyone)
}
