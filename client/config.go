package client

import (
	"os"
	"time"
)

//Config for Client
type Config struct {
	ServiceName string
	URL         string
	SyncTime    time.Duration
}

//NewConfig struct with defaults
func NewConfig() *Config {
	return &Config{
		ServiceName: defaultServiceName(),
		URL:         defaultConfigURL(),
		SyncTime:    defaultSyncTime(),
	}
}

func defaultServiceName() string {
	serviceName, ok := os.LookupEnv("SERVICE_NAME")
	if ok {
		return serviceName
	}
	return ""
}

func defaultConfigURL() string {
	return "http://localhost:7654"
}

func defaultSyncTime() time.Duration {
	return time.Second * 30
}
