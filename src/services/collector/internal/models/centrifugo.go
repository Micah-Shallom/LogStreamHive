package models

import (
	"log"

	"github.com/centrifugal/gocent"
)

type CentrifugoConfig struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
}

type CentrifugoClient struct {
	Client *gocent.Client
	Logger *log.Logger
}
