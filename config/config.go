package config

import (
	"clip/cmd"
	"fmt"
	"sync"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	App        App
	PostgreSQL PostgreSQL
	Log        Log
}

type App struct {
	Host      string
	Port      int
	Debug     bool
	Whitelist []string
}

type Log struct {
	Prefix     string
	Dir        string
	LevelDebug bool
}

type PostgreSQL struct {
	Username     string
	Password     string
	Host         string
	Port         int
	Db           string
	Debug        bool
	MaxIdleConns int
	MaxOpenConns int
}

func init() {
	// Init CLI commands
	cmd.Root().Use = "bin/clip --config <Config path>"
	cmd.Root().Short = "clip - Provide API for clip"
	cmd.Root().Long = "clip"

	cmd.SetRunFunc(load)
}

func load() {
	once.Do(func() {
		if err := cmd.GetViper().Unmarshal(&conf); err != nil {
			fmt.Println("load viper fail")
		}
	})
}

func Load() {
	load()
}

func Get() Config {
	load()
	return conf
}

func GetPostgreSQL() PostgreSQL { return conf.GetPostgreSQL() }
func (c Config) GetPostgreSQL() PostgreSQL {
	return c.PostgreSQL
}
