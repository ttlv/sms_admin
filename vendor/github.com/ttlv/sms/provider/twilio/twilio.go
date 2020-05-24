package twilio

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/sfreiberg/gotwilio"
	"github.com/tidwall/gjson"
	"github.com/ttlv/sms"
	"github.com/ttlv/sms/config"
)

type TwilioProvider struct {
	DB *gorm.DB
}

func New() TwilioProvider {
	cfg := config.MustGetConfig()
	db, err := gorm.Open("mysql", cfg.DB)
	if err != nil {
		panic(err)
	}
	return TwilioProvider{DB: db}
}

func (provider TwilioProvider) GetCode() string {
	return "Twilio"
}

func (provider TwilioProvider) AvailableCountries() []string {
	return []string{}
}

func (provider TwilioProvider) Available(s *sms.SendParams) bool {
	b := sms.SmsBrand{}
	provider.DB.Where("name = ?", s.Brand).First(&b)
	if !b.EnableTwilio || b.TwilioAccountsID == "" || b.TwilioAuthToken == "" || b.TwilioSendNumber == "" {
		return false
	}
	return true
}

func (provider TwilioProvider) Send(params sms.SendParams) (string, string, error) {
	c := config.MustGetConfig()
	brand := sms.SmsBrand{}
	provider.DB.Where("name = ?", params.Brand).First(&brand)
	twilioClient := gotwilio.NewTwilioClient(brand.TwilioAccountsID, brand.TwilioAuthToken)
	message, _, err := twilioClient.SendSMS(brand.TwilioSendNumber, params.Phone, params.Content, c.TwilioCallBack, "")
	if err != nil {
		return "", "", err
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		return "", "", err
	}
	return string(jsonData), gjson.Get(string(jsonData), "sid").String(), err
}
