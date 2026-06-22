package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	once sync.Once
	cfg  *config
)

type config struct {
	debug bool
	pg    Postgres
}

func (c *config) Debug() bool        { return c.debug }
func (c *config) Postgres() Postgres { return c.pg }

func Config() *config {
	once.Do(func() {
		viper.AutomaticEnv()

		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AddConfigPath("../..")

		_ = viper.ReadInConfig()

		cfg = &config{
			debug: viper.GetBool("DEBUG"),
			pg: Postgres{
				Host:     viper.GetString("POSTGRES_HOST"),
				Port:     viper.GetString("POSTGRES_PORT"),
				User:     viper.GetString("POSTGRES_USER"),
				Password: viper.GetString("POSTGRES_PASSWORD"),
				DBName:   viper.GetString("POSTGRES_DBNAME"),
				Schema:   viper.GetString("POSTGRES_SCHEMA"),
			},
		}
	})
	return cfg
}
