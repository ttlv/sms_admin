package amazon

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ttacon/libphonenumber"
	"github.com/ttlv/sms"
	"github.com/ttlv/sms/config"
)

type AmazonProvider struct {
	DB *gorm.DB
}

type Message struct {
	Mobile             string `json:"mobile"`
	Content            string `json:"content"`
	RequestTime        int64  `json:"requestTime"`
	RequestValidPeriod int    `json:"requestValidPeriod"`
}

func New() AmazonProvider {
	cfg := config.MustGetConfig()
	db, err := gorm.Open("mysql", cfg.DB)
	if err != nil {
		panic(err)
	}
	return AmazonProvider{DB: db}
}

func (provider AmazonProvider) GetCode() string {
	return "Amazon"
}

func (provider AmazonProvider) AvailableCountries() []string {
	return []string{}
}

func (provider AmazonProvider) Available(s *sms.SendParams) bool {
	b := sms.SmsBrand{}
	provider.DB.Where("name = ?", s.Brand).First(&b)
	if !b.EnableAWS || b.AWSRegion == "" || b.AWSAccessKeyID == "" || b.AWSSecretAccessKey == "" {
		return false
	}
	return true
}

func (provider AmazonProvider) Send(params sms.SendParams) (string, string, error) {
	brand := sms.SmsBrand{}
	phonenumber, _ := libphonenumber.Parse(params.Phone, params.Country)
	params.Phone = strings.Replace(libphonenumber.Format(phonenumber, libphonenumber.E164), " ", "", -1)
	provider.DB.Where("name = ?", params.Brand).First(&brand)
	creds := credentials.NewStaticCredentials(brand.AWSAccessKeyID, brand.AWSSecretAccessKey, "")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: creds,
			Region:      aws.String(brand.AWSRegion),
		},
	}))
	svc := sns.New(sess)
	param := &sns.PublishInput{
		Message:     aws.String(params.Content),
		PhoneNumber: aws.String(params.Phone),
	}
	resp, err := svc.Publish(param)
	if err != nil {
		return "", "", err
	}
	return *resp.MessageId, *resp.MessageId, err
}
