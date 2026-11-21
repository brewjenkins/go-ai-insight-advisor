package config

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Env                string   `mapstructure:"env"`
	Debug              bool     `mapstructure:"debug"`
	CorsAllowedOrigins []string `mapstructure:"cors_allowed_origins"`
}

type Config struct {
	App AppConfig `mapstructure:"app"`
}

func NewConfig() *Config {
	return &Config{
		App: AppConfig{
			Server: struct {
				Host string `mapstructure:"host"`
				Port string `mapstructure:"port"`
			}{
				Host: "localhost",
				Port: "8080",
			},
			Env:                "development",
			Debug:              true,
			CorsAllowedOrigins: []string{"*"},
		},
	}
}
