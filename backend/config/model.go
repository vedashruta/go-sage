package config

import "time"

var (
	Config config
)

type config struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Timeout      string        `json:"timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	Storage      string
}
