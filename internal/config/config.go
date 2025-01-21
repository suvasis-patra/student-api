package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-required:"true"`
	DbPath string `yaml:"db_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "finds the config variable")
		flag.Parse()
		configPath = *flags
		fmt.Println(configPath,*flags)
		if configPath == "" {
			log.Fatal("Configuration path is not set!")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file not found! %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cann't read config file: %s", err.Error())
	}
	return &cfg
}
