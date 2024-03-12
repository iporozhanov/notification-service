package config

import "time"

type Config struct {
	HTTP                     HTTP              `yaml:"http"`
	DB                       DB                `yaml:"db"`
	Senders                  map[string]string `yaml:"senders"`
	Mailgun                  Mailgun           `yaml:"mailgun"`
	Twilio                   Twilio            `yaml:"twilio"`
	NotificationMaxAttempts  int64             `yaml:"notification_max_attempts"`
	NotificationListenTicker time.Duration     `yaml:"notification_listen_ticker"`
	Auth                     Auth              `yaml:"auth"`
	Sinch                    Sinch             `yaml:"sinch"`
	Slack                    Slack             `yaml:"slack"`
}

type HTTP struct {
	Port                 string        `yaml:"port"`
	RatelimitTimeout     time.Duration `yaml:"ratelimit_timeout"`
	RatelimitMaxRequests int64         `yaml:"ratelimit_max_requests"`
}

type DB struct {
	URL string `yaml:"url"`
}

type Mailgun struct {
	From       string `yaml:"from"`
	Domain     string `yaml:"domain"`
	PrivateKey string `yaml:"private_key"`
}

type Twilio struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type Sinch struct {
	From          string `yaml:"from"`
	ServicePlanID string `yaml:"service_plan_id"`
	APIKey        string `yaml:"api_key"`
}

type Auth struct {
	APIKey string `yaml:"api_key"`
}

type Slack struct {
	Token string `yaml:"token"`
}

func (c *Config) Default() {
	c.HTTP.Port = "8080"
	c.HTTP.RatelimitTimeout = 10 * time.Second
	c.HTTP.RatelimitMaxRequests = 10
	c.DB.URL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	c.NotificationMaxAttempts = 3
	c.NotificationListenTicker = 10 * time.Second
}
