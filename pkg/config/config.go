package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	All        = config{}
	Global     = &All.Global
	HttpServer = &All.HttpServer
	Storage    = &All.Storage
	Log        = &All.Log
	Sentry     = &All.Sentry
)

func init() {
	if !All.parse() {
		showHelp()
		os.Exit(1)
	}
}

type global struct {
	CertFile string `mapstructure:"certFile"`
	KeyFile  string `mapstructure:"keyFile"`
}

type httpServer struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
}

type storage struct {
	Redis redis `mapstructure:"redis"`
	Mongo mongo `mapstructure:"mongo"`
}

type redis struct {
	Host     string                 `mapstructure:"host"`
	Port     int                    `mapstructure:"port"`
	Database int                    `mapstructure:"database"`
	Password string                 `mapstructure:"password"`
	Extra    map[string]interface{} `mapstructure:"extra"`
}

type mongo struct {
	Host     string                 `mapstructure:"host"`
	Port     int                    `mapstructure:"port"`
	Database string                 `mapstructure:"database"`
	User     string                 `mapstructure:"user"`
	Password string                 `mapstructure:"password"`
	Timeout  uint8                  `mapstructure:"timeout"`
	Extra    map[string]interface{} `mapstructure:"extra"`
}

type log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type sentry struct {
	Url     string `mapstructure:"url"`
	Debug   bool   `mapstructure:"debug"`
	Timeout uint8  `mapstructure:"timeout"`
}

type config struct {
	Global     global     `mapstructure:"global"`
	HttpServer httpServer `mapstructure:"httpServer"`
	Storage    storage    `mapstructure:"storage"`
	Log        log        `mapstructure:"log"`
	Sentry     sentry     `mapstructure:"sentry"`
	CfgFile    string
}

func showHelp() {
	fmt.Printf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {configs file}")
	fmt.Println("      -h (show help info)")
}

func (c *config) load() bool {
	_, err := os.Stat(c.CfgFile)

	if err != nil {
		return false
	}

	viper.SetConfigFile(c.CfgFile)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()

	if err != nil {
		fmt.Printf("configs file %s read failed. %v\n")
		return false
	}

	err = viper.GetViper().UnmarshalExact(c)

	if err != nil {
		fmt.Printf("configs file %s loaded failed. %v\n", c.CfgFile, err)
		return false
	}

	return true
}

func (c *config) parse() bool {
	flag.StringVar(&c.CfgFile, "c", "configs/config.toml", "configs file")

	help := flag.Bool("h", false, "help info")

	flag.Parse()

	if !c.load() {
		return false
	}

	if *help {
		showHelp()
		return false
	}

	return true
}
