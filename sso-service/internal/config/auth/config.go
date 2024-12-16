package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string     `yaml:"env" env-default:"local"`
	StoragePath   string     `yaml:"storage_path" env-required:"true"`
	GRPC          GRPCConfig `yaml:"grpc"`
	MigrationPath string
	TokenTTL      time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	serviceName := "auth"
	configPath := fetchConfigPath(serviceName)
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(configPath)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath(serviceName string) string {
	var res string

	flag.StringVar(&res, "config", fmt.Sprintf("../../config/config_%s.yaml", serviceName), "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	fmt.Println("Config path from flag:", res)
	fmt.Println("Config path from env:", os.Getenv("CONFIG_PATH"))

	return res
}
