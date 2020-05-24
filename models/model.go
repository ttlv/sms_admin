package models

import (
	"time"
)

var (
	smsState = [][]string{{"0", "发送中"}, {"1", "已发送"}, {"2", "发送失败"}, {"3", "用户已收到"}}
)

type SendParamForm struct {
	Phone   string
	Content string
}

type SmsRecord struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Brand     string
	// 客户端发送的原始请求参数
	RawParam       string
	Phone          string
	State          int8
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
type SmsAvailable struct {
	ID              uint `gorm:"primary_key"`
	SmsBranID       uint
	SmsBrand        SmsBrand `gorm:"foreignkey:SmsBranID"`
	AvailableAmount int64
	Note            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
