package config

import "flag"

type AppConfig struct {
	Port        int    // Server addr
	DatabaseURL string // Database connection URL
}

// FromFlags creates server config based on passed arguments.
func FromFlags() *AppConfig {
	c := &AppConfig{}

	flag.IntVar(&c.Port, "addr", 8080, "Port on which server should start")
	flag.StringVar(&c.DatabaseURL, "db-url", "postgresql://localhost:5432/shp", "Database connection URL")

	return c
}
