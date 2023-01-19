package utility

import "github.com/spf13/viper"

// Config file stores configuration of application

// Values read by viper viper from a config file or enviroment variable

type Config struct {
	SeverAddress string `mapstructure:"SEVER_ADDRESS"`
	FromEmail    string `mapstructure:"FROM_EMAIL"`
	Password     string `mapstructure:"PASSWORD"`
	Host         string `mapstructure:"HOST"`
	MailPort     string `mapstructure:"MAIL_PORT"`
	GrpcPort     string `mapstructure:"GRPC_PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
