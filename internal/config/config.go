package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name        string
		Port        int
		Environment string
	}
	MongoDB struct {
		URI        string `mapstructure:"uri"`
		Database   string `mapstructure:"database"`
		Collection string `mapstructure:"collection"`
	} `mapstructure:"mongodb"`
	Elasticsearch struct {
		Addresses []string `mapstructure:"addresses"`
		Username  string   `mapstructure:"username"`
		Password  string   `mapstructure:"password"`
		Index     string   `mapstructure:"index"`
	} `mapstructure:"elasticsearch"`
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		GroupID string   `mapstructure:"group_id"`
		Topic   struct {
			ProductCreated  string `mapstructure:"product_created"`
			ProductUpdated  string `mapstructure:"product_updated"`
			ProductDeleted  string `mapstructure:"product_deleted"`
			ProductViewsInc string `mapstructure:"product_views_inc"`
			ProductBuysInc  string `mapstructure:"product_buys_inc"`
		} `mapstructure:"topic"`
	} `mapstructure:"kafka"`
	Logging struct {
		Level  string
		Format string
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
