package config

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

var Conf *Config

// multi layers for yaml file
type Config struct {
	Server   *Server             `yaml:"server"`
	MySQL    *MySQL              `yaml:"mysql"`
	Redis    *Redis              `yaml:"redis"`
	Etcd     *Etcd               `yaml:"etcd"`
	Services map[string]*Service `yaml:"services"`
	Domain   map[string]*Domain  `yaml:"domain"`
}

// server: http service
type Server struct {
	Port      string `yaml:"port"`
	Version   string `yaml:"version"`
	JwtSecret string `yaml:"jwtSecret"`
}

type MySQL struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
	Charset    string `yaml:"charset"`
}

type Redis struct {
	UserName string `yaml:"userName"`
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

// item -> [item]
type Service struct {
	Name        string `yaml:"name"`
	LoadBalance bool   `yaml:"loadBalance"`
	// using slice
	Addr []string `yaml:"addr"`
}

type Domain struct {
	Name string `yaml:"name"`
}

// set viper config
func InitConfig() {
	workdir, _ := os.Getwd()
	configpath := path.Join(workdir, "config")

	// set config.yaml for viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetConfigType(configpath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	//unmarshal config using yaml struct
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}
