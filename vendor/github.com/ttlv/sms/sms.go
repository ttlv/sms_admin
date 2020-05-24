package sms

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	RecordState_Sending   RecordState = 0
	RecordState_Success   RecordState = 1
	RecordState_Failure   RecordState = 2
	RecordState_Delivered RecordState = 3
)

type SmsProvider interface {
	GetCode() string
	// 空数组代表支持所有国家的发送
	AvailableCountries() []string
	Available(*SendParams) bool
	Send(SendParams) (resp string, externalID string, err error)
}

type PublishData struct {
	SmsRecordId   uint
	SentProviders []string
	SendParams    *SendParams
}

type SmsQueue interface {
	Publish(data *PublishData)
	Liveness() bool
}

type RecordState int32

type SmsRecord struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Brand     string
	// 客户端发送的原始请求参数
	RawParam       string
	Phone          string
	State          RecordState
	Error          string `gorm:"size:1024"`
	ProviderResp   string `gorm:"size:1024"`
	ExternalID     string
	Sender         string
	LastSendAt     *time.Time
	LastCallbackAt *time.Time
}

type SmsFailureRecord struct {
	ID           uint `gorm:"primary_key" json:"-"`
	SmsRecordId  uint
	SmsRecord    SmsRecord
	CreatedAt    *time.Time `json:"-"`
	UpdatedAt    *time.Time `json:"-"`
	ProviderName string
	Phone        string
	Error        string `gorm:"size:1024"`
}

type SmsBrand struct {
	ID                 uint `gorm:"primary_key"`
	Name               string
	Token              string
	TwilioAccountsID   string
	TwilioAuthToken    string
	TwilioSendNumber   string
	YunPianAppKey      string
	YunPianHost        string
	EmayAppID          string
	EmayAppKey         string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	ChuangLanAccount   string
	ChuangLanPassword  string
	EnableAWS          bool
	EnableTwilio       bool
	EnableYunPian      bool
	EnableEmay         bool
	EnableChuangLan    bool
}

type SmsSetting struct {
	ID            uint `gorm"primary_key"`
	Content       string
	ProviderSorts string
}

type ChuangLanCallBack struct {
	Clcode string `json:"clcode"`
}

type SmsAvailable struct {
	ID              uint `gorm:"primary_key"`
	SmsBranID       uint
	SmsBrand        SmsBrand `gorm:"foreignkey:SmsBranID"`
	AvailableAmount int64
	Note            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ApiParams struct {
	Brand   string
	Phone   string
	Content string
}

type HttpParam struct {
	Brand   string `json:"brand"`
	Phone   string `json:"phone"`
	Content string `json:"content"`
}

func (record SmsRecord) Content() string {
	sendParam := SendParams{}
	json.Unmarshal([]byte(record.RawParam), &sendParam)
	return sendParam.Content
}

func FormattedProviderSorts(db *gorm.DB) []string {
	setting := &SmsSetting{}
	db.First(&setting)
	return strings.Split(strings.Replace(setting.ProviderSorts, " ", "", -1), ",")
}
