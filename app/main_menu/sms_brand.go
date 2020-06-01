package main_menu

import (
	"github.com/qor/admin"
	"github.com/ttlv/sms/provider/amazon"
	"github.com/ttlv/sms/provider/chuanglan"
	"github.com/ttlv/sms/provider/emay"
	"github.com/ttlv/sms/provider/twilio"
	"github.com/ttlv/sms/provider/yunpian"
	"github.com/ttlv/sms_admin/config/application"
	"github.com/ttlv/sms_admin/models"
	"github.com/ttlv/sms_admin/utils"
)

func ConfigBrandRes(application *application.Application) {
	res := application.Admin.AddResource(&models.SmsBrand{})
	res.Action(&admin.Action{
		Name:  "TestAmazon",
		Label: "测试亚马逊",
		Handler: func(actionArgument *admin.ActionArgument) error {
			params := utils.ConstructParam(actionArgument)
			resp, externalID, err := amazon.New().Send(*params)
			if err != nil {
				return err
			}
			utils.SaveSmsRecord(actionArgument.Context.DB, params, "Amazon", resp, externalID)
			return nil
		},
		Resource: application.Admin.NewResource(models.SendParamForm{}),
		Modes:    []string{"edit"},
	})
	res.Action(&admin.Action{
		Name:  "TestEmay",
		Label: "测试亿美",
		Handler: func(actionArgument *admin.ActionArgument) error {
			params := utils.ConstructParam(actionArgument)
			resp, externalID, err := emay.New().Send(*params)
			if err != nil {
				return err
			}
			utils.SaveSmsRecord(actionArgument.Context.DB, params, "Emay", resp, externalID)
			return nil
		},
		Resource: application.Admin.NewResource(models.SendParamForm{}),
		Modes:    []string{"edit"},
	})
	res.Action(&admin.Action{
		Name:  "Test twillp",
		Label: "测试twilio",
		Handler: func(actionArgument *admin.ActionArgument) error {
			params := utils.ConstructParam(actionArgument)
			resp, externalID, err := twilio.New().Send(*params)
			if err != nil {
				return err
			}
			utils.SaveSmsRecord(actionArgument.Context.DB, params, "Twilio", resp, externalID)
			return nil
		},
		Resource: application.Admin.NewResource(models.SendParamForm{}),
		Modes:    []string{"edit"},
	})
	res.Action(&admin.Action{
		Name:  "Test yunpian",
		Label: "测试云片",
		Handler: func(actionArgument *admin.ActionArgument) error {
			params := utils.ConstructParam(actionArgument)
			resp, externalID, err := yunpian.New().Send(*params)
			if err != nil {
				return err
			}
			utils.SaveSmsRecord(actionArgument.Context.DB, params, "YunPian", resp, externalID)
			return nil
		},
		Resource: application.Admin.NewResource(models.SendParamForm{}),
		Modes:    []string{"edit"},
	})
	res.Action(&admin.Action{
		Name:  "Test ChuangLan",
		Label: "测试创蓝",
		Handler: func(argument *admin.ActionArgument) error {
			params := utils.ConstructParam(argument)
			resp, externalID, err := chuanglan.New().Send(*params)
			if err != nil {
				return err
			}
			utils.SaveSmsRecord(argument.Context.DB, params, "ChuangLan", resp, externalID)
			return nil
		},
		Resource: application.Admin.NewResource(models.SendParamForm{}),
		Modes:    []string{"edit"},
	})

	res.IndexAttrs("Name", "Token", "TwilioAccountsID", "TwilioAuthToken", "TwilioSendNumber", "YunPianAppKey", "EmayAppID", "EmayAppKey", "AWSAccessKeyID", "AWSSecretAccessKey", "AWSRegion", "ChuangLanAccount", "ChuangLanPassword")
	res.NewAttrs(&admin.Section{
		Title: "Base Config",
		Rows: [][]string{
			{"Name", "Token"},
		},
	},
		&admin.Section{
			Title: "Twilio Config",
			Rows: [][]string{
				{"TwilioAccountsID", "TwilioAuthToken"},
				{"TwilioSendNumber", "EnableTwilio"},
			},
		},
		&admin.Section{
			Title: "YunPian Config",
			Rows: [][]string{
				{"YunPianAppKey", "YunPianHost"},
				{"EnableYunPian"},
			},
		},
		&admin.Section{
			Title: "Emay Config",
			Rows: [][]string{
				{"EmayAppID", "EmayAppKey"},
				{"EnableEmay"},
			},
		},
		&admin.Section{
			Title: "Aws Config",
			Rows: [][]string{
				{"AWSAccessKeyID", "AWSSecretAccessKey"},
				{"AWSRegion", "EnableAWS"},
			},
		},
		&admin.Section{
			Title: "ChuangLan",
			Rows: [][]string{
				{"ChuangLanAccount", "ChuangLanPassword"},
				{"EnableChuangLan"},
			},
		})
	res.EditAttrs(res.NewAttrs())
}
