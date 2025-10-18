package models

type Config struct {
	LogFiles      []string         `yaml:"log_files" json:"log_files"`
	CheckInterval float64          `yaml:"check_interval" json:"check_interval"`
	Centrifugo    CentrifugoConfig `yaml:"centrifugo" json:"centrifugo"`
	ChannelID     string           `yaml:"channel_id" json:"channel_id"`
	Nats          NatsConfig       `yaml:"nats"`
}
