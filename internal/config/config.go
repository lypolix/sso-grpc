package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port   int           `yaml:"port"`
	Timout time.Duration `yaml:"timeout"`
}


var configPath string

func init() {

	flag.StringVar(&configPath, "config", "", "path to config file")
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("empty config path")
	}
	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if path == "" {
		panic("empty config path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist:" + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}
	return configPath
}