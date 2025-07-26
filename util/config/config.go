package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver         string `mapstructure:"DB_DRIVER"`
	MigrationUrl     string `mapstructure:"MIGRATION_URL"`
	DBDSource        string `mapstructure:"DB_SOURCE"`
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
	RedisAddress     string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword    string `mapstructure:"REDIS_PASSWORD"`
	RedisDB          int    `mapstructure:"REDIS_DB"`
	SecretKey        string `mapstructure:"SECRET_KEY"`
	MaxMsgBuffSize   int    `mapstructure:"MAX_MESSAGE_BUFFER_SIZE"`
	MaxEntryBuffSize int    `mapstructure:"MAX_ENTRY_BUFFER_SIZE"`
	LoggerLevel      int8   `mapstructure:"LOGGER_LEVEL"`
}

func InitConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return config, err

}
