package gofast

import "time"

type Config struct {
	// fasthttp client configurations
	Name                     string
	NoDefaultUserAgentHeader bool
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration

	// RequestEncoder encode request before send
	RequestEncoder RequestEncoder

	// ResponseDecoder decode response after send
	ResponseDecoder ResponseDecoder

	// ErrorHandler handle the status code without 2xx
	ErrorHandler ErrorHandler
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Name:                     "gofast",
	NoDefaultUserAgentHeader: false,
	ReadTimeout:              6 * time.Second,
	WriteTimeout:             6 * time.Second,
	RequestEncoder:           JSONEncoder,
	ResponseDecoder:          JSONDecoder,
	ErrorHandler:             defaultErrorHandler,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Name == "" {
		cfg.Name = ConfigDefault.Name
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = ConfigDefault.ReadTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = ConfigDefault.WriteTimeout
	}
	if cfg.RequestEncoder == nil {
		cfg.RequestEncoder = ConfigDefault.RequestEncoder
	}
	if cfg.ResponseDecoder == nil {
		cfg.ResponseDecoder = ConfigDefault.ResponseDecoder
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = ConfigDefault.ErrorHandler
	}
	return cfg
}
