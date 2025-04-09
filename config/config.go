package config

import (
	"log"

	"github.com/spf13/viper"
)

func ConfigInit() {
	viper.SetConfigFile("config/config.yml") // 指定配置文件路径
	// viper.SetConfigName("config")
	// viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("ReadInConfig Error: ", err)
		return
	}
}
