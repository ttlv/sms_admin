package chuanglan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"unsafe"

	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
	"github.com/ttlv/sms"
	"github.com/ttlv/sms/config"
	"strings"
)

type ChuangLanProvider struct {
	DB *gorm.DB
}

func New() ChuangLanProvider {
	cfg := config.MustGetConfig()
	db, err := gorm.Open("mysql", cfg.DB)
	if err != nil {
		panic(err)
	}
	return ChuangLanProvider{DB: db}
}

func (provider ChuangLanProvider) GetCode() string {
	return "ChuangLan"
}

func (provider ChuangLanProvider) AvailableCountries() []string {
	return []string{"CN"}
}

func (provider ChuangLanProvider) Available(s *sms.SendParams) bool {
	b := sms.SmsBrand{}
	provider.DB.Where("name = ?", s.Brand).First(&b)
	if b.ChuangLanAccount == "" || b.ChuangLanPassword == "" {
		return false
	}
	return true
}

func (provider ChuangLanProvider) Send(params sms.SendParams) (string, string, error) {
	var (
		b = sms.SmsBrand{}
	)
	apiParams := make(map[string]interface{})
	provider.DB.Where("name = ?", params.Brand).First(&b)
	//请登录zz.253.com获取API账号、密码以及短信发送的URL
	apiParams["account"] = b.ChuangLanAccount                           //创蓝API账号
	apiParams["password"] = b.ChuangLanPassword                         //创蓝API密码
	apiParams["phone"] = strings.Replace(params.Phone, " ", "", -1)[3:] //手机号码
	//设置您要发送的内容：其中“【】”中括号为运营商签名符号，多签名内容前置添加提交
	apiParams["msg"] = url.QueryEscape(params.Content)
	apiParams["report"] = "true"
	bytesData, err := json.Marshal(apiParams)
	if err != nil {
		return "", "", err
	}
	reader := bytes.NewReader(bytesData)
	url := "https://smssh1.253.com/msg/send/json" //短信发送URL
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", "", err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	body := (*string)(unsafe.Pointer(&respBytes))
	if gjson.Get(*body, "code").String() == "0" {
		return *body, gjson.Get(*body, "msgId").String(), nil
	}
	return "", "", fmt.Errorf(gjson.Get(*body, "errorMsg").String())
}
