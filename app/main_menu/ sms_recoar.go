package main_menu

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/ttlv/sms_admin/config/application"
	"github.com/ttlv/sms_admin/models"
	"github.com/ttlv/sms_admin/utils"
	"html/template"
)

func ConfigSmsRecoardRes(application *application.Application) {
	res := application.Admin.AddResource(&models.SmsRecord{})
	res.UseTheme("readonly")
	//Brand Filter
	res.Filter(&admin.Filter{
		Name:  "Brand",
		Label: "服务名",
		Config: &admin.SelectOneConfig{
			Collection: func(i interface{}, c *qor.Context) (result [][]string) {
				brands := []models.SmsBrand{}
				c.DB.Find(&brands)
				for _, b := range brands {
					result = append(result, []string{b.Name, b.Name})
				}
				return
			},
		},
	})

	//State Filter
	res.Filter(&admin.Filter{
		Name:  "State",
		Label: "状态",
		Config: &admin.SelectOneConfig{
			Collection: func(i interface{}, c *qor.Context) (result [][]string) {
				return [][]string{{"0", "Sending"}, {"1", "Success"}, {"2", "Failure"}, {"3", "Deliverd"}}
			},
		},
	})

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
		Name:  "Brand",
		Label: "服务名",
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			return record.(*models.SmsRecord).Brand
		},
	})

	res.Meta(&admin.Meta{
		Name:  "State",
		Label: "状态",
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			if record.(*models.SmsRecord).State == 0 {
				return template.HTML("<div style='display:block;border-radius:50%;height:10px;width:10px;background-color:yellow'>️</div>")
			} else if record.(*models.SmsRecord).State == 1 {
				return template.HTML("<div class='delivered' style='width:10px;height:10px;border-radius:50%;background:yellow;background-image: linear-gradient(to right, transparent 50%, green 0)'>️</div>")
			} else if record.(*models.SmsRecord).State == 2 {
				return template.HTML("<div style='width:10px;height:10px;background:red;border-radius:50%'>️</div>")
			} else if record.(*models.SmsRecord).State == 3 {
				return template.HTML("<div style='width:10px;height:10px;background:#09e609c2;border-radius:50%'>️</div>")
			}
			return
		},
	})

	res.Meta(&admin.Meta{
		Name: "TimeInterval",
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			if record.(*models.SmsRecord).LastCallbackAt != nil && record.(*models.SmsRecord).LastSendAt != nil {
				return record.(*models.SmsRecord).LastCallbackAt.Unix() - record.(*models.SmsRecord).LastSendAt.Unix()
			}
			return "等待回调中"
		},
	})

	res.Meta(&admin.Meta{
		Name:  "CreatedAt",
		Label: "创建时间",
	})

	res.Meta(&admin.Meta{
		Name:  "Phone",
		Label: "手机号码",
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			return record.(*models.SmsRecord).Phone
		},
	})

	res.Meta(&admin.Meta{
		Name:  "Sender",
		Label: "短信服务商",
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			return record.(*models.SmsRecord).Sender
		},
	})
	res.Meta(&admin.Meta{
		Name:  "Error",
		Label: "错误",
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			return record.(*models.SmsRecord).Error
		},
	})

	res.IndexAttrs("CreatedAt", "Brand", "Phone", "Sender", "State", "Error", "TimeInterval")
	res.Permission = roles.Deny(roles.Create, roles.Anyone).Deny(roles.Update, roles.Anyone).Deny(roles.Delete, roles.Anyone)
}
