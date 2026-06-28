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
	grpc  Grpc
}

func (c *config) Debug() bool        { return c.debug }
func (c *config) Postgres() Postgres { return c.pg }
func (c *config) Grpc() Grpc         { return c.grpc }

func Config() *config {
	once.Do(func() {
		viper.AutomaticEnv()

		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AddConfigPath("../..")

		_ = viper.ReadInConfig()

		grpcPort := viper.GetInt("GRPC_PORT")
		if grpcPort == 0 {
			grpcPort = 50051
		}

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
			grpc: Grpc{
				Port: grpcPort,
			},
		}
	})
	return cfg
}
