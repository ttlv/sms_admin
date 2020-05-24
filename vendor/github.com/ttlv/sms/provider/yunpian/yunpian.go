package yunpian

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"github.com/ttacon/libphonenumber"
	"github.com/ttlv/sms"
	"github.com/ttlv/sms/config"
)

type YunPianProvider struct {
	DB *gorm.DB
}

func New() YunPianProvider {
	cfg := config.MustGetConfig()
	db, err := gorm.Open("mysql", cfg.DB)
	if err != nil {
		panic(err)
	}
	return YunPianProvider{DB: db}
}

func (provider YunPianProvider) GetCode() string {
	return "YunPian"
}

func (provider YunPianProvider) AvailableCountries() []string {
	return []string{"CN"}
}

func (provider YunPianProvider) Available(s *sms.SendParams) bool {
	b := sms.SmsBrand{}
	provider.DB.Where("name = ?", s.Brand).First(&b)
	if !b.EnableYunPian || b.YunPianAppKey == "" || b.YunPianHost == "" {
		return false
	}
	return true
}

func (provider YunPianProvider) Send(params sms.SendParams) (string, string, error) {
	c := config.MustGetConfig()
	phonenumber, _ := libphonenumber.Parse(params.Phone, params.Country)
	params.Phone = strings.Replace(libphonenumber.Format(phonenumber, libphonenumber.E164), " ", "", -1)

	brand := sms.SmsBrand{}
	host := "https://sms.yunpian.com"
	provider.DB.Where("name = ?", params.Brand).First(&brand)
	if brand.YunPianHost != "" {
		host = brand.YunPianHost
	}

	data := url.Values{"apikey": {brand.YunPianAppKey}, "mobile": {params.Phone}, "text": {params.Content}, "callback_url": {c.YunPianCallBack}}
	resp, err := http.PostForm(host+"/v1/sms/send.json", data)
	if err != nil {
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	code := gjson.Get(string(body), "code").Int()
	if code == 0 || code == 8 {
		return string(body), gjson.Get(string(body), "result.sid").String(), nil
	}
	return "", "", fmt.Errorf(string(body))
}
