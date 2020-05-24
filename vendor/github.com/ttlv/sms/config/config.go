package config

import "github.com/jinzhu/configor"

type Config struct {
	DB              string
	Address         string
	AMQPDial        string
	APITOKEN        string
	TwilioCallBack  string
	YunPianCallBack string
	APIKEY          string
	Port            string
}

var _config *Config

func MustGetConfig() Config {
	if _config != nil {
		return *_config
	}

	_config = &Config{}
	err := configor.New(&configor.Config{ENVPrefix: "SMS"}).Load(_config)
	if err != nil {
		panic(err)
	}

	return *_config
}

