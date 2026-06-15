package config

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	DB_URL             string `json:"DB_URL"`
	REDIS_ADDR         string
	URL_TTL            int
	SHORT_CODE_MIN_LEN uint8
	BASE_URL           string
	API_PORT           int
	ChainCertPath      string
	PemCertPath        string
}

func LoadConfig(envType string) *Config {
	conf := Config{}
	v := viper.New()
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error while fetching user's Home Dir")
	}
	configDir := path.Join(dir, "/.shortner")
	v.SetConfigType("json")
	v.SetConfigFile(fmt.Sprintf("%s/config_%s.json", configDir, envType))

	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			slog.Warn("couldn't find Config file: reading env")
			v.AutomaticEnv()
		} else {
			log.Fatal("Unable to read config file: ", err)
		}
	}

	slog.Info(fmt.Sprintf("read config file %s", v.ConfigFileUsed()))

	if v.IsSet("DB_URL") {
		conf.DB_URL = v.GetString("DB_URL")
	}
	if v.IsSet("BASE_URL") {
		conf.BASE_URL = v.GetString("BASE_URL")
	}
	if v.IsSet("API_PORT") {
		conf.API_PORT = v.GetInt("API_PORT")
	}
	if v.IsSet("CHAIN_PATH") {
		conf.ChainCertPath = v.GetString("CHAIN_PATH")
	}
	if v.IsSet("PEM_PATH") {
		conf.PemCertPath = v.GetString("PEM_PATH")
	}
	if v.IsSet("REDIS_ADDR") {
		conf.REDIS_ADDR = v.GetString("REDIS_ADDR")
	}
	if v.IsSet("URL_TTL") {
		conf.URL_TTL = v.GetInt("URL_TTL")
	} else {
		// default ttl for 60 Minutes
		conf.URL_TTL = 60
	}
	if v.IsSet("SHORT_CODE_MIN_LEN") {
		conf.SHORT_CODE_MIN_LEN = v.GetUint8("SHORT_CODE_MIN_LEN")
	} else {
		// default ttl for 60 Minutes
		conf.SHORT_CODE_MIN_LEN = uint8(6)
	}

	if conf.DB_URL == "" {
		log.Fatal("DB_URL is required")
	}
	if conf.REDIS_ADDR == "" {
		log.Fatal("REDIS_ADDR is required")
	}
	if conf.API_PORT == 0 {
		log.Fatal("API_PORT is required")
	}
	if conf.BASE_URL == "" {
		log.Fatal("BASE_URL is required")
	}

	return &conf
}
