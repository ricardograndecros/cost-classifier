// cost-classifier/backend/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	SecretKey         string `json:"secret_key"`
	NordigenSecretId  string `json:"nordigen_secret_id"`
	NordigenSecretKey string `json:"nordigen_secret_key"`
	BankID            string `json:"bank_id"`
	EndUserID         string `json:"end_user_id"`
}

var AppConfig = Config{}

func InitConfig() {
	viper.SetConfigName(".config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return
	}

	AppConfig.SecretKey = viper.GetString("secret_key")
	AppConfig.NordigenSecretId = viper.GetString("nordigen_secret_id")
	AppConfig.NordigenSecretKey = viper.GetString("nordigen_secret_key")
	AppConfig.BankID = viper.GetString("bank_id")
	AppConfig.EndUserID = viper.GetString("end_user_id")
}
