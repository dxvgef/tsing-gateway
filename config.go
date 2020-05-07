package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// LocalConfig 全局配置
var LocalConfig struct {
	Listener struct {
		IP                string        `yaml:"ip"`
		HTTPPort          int           `yaml:"httpPort"`
		HTTPSPort         int           `yaml:"httpsPort"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
		QuitWaitTimeout   time.Duration `yaml:"quitWaitTimeout"`
		HTTP2             bool          `yaml:"http2"`
		Debug             bool          `yaml:"debug"`
	} `yaml:"listener"`
	Service struct {
		Name       string `yaml:"name"`
		CenterAddr string `yaml:"centerAddr"`
	} `yaml:"service"`
	Logger struct {
		Level      string      `yaml:"level"`
		FilePath   string      `yaml:"filePath"`
		FileMode   os.FileMode `yaml:"fileMode"`
		Encode     string      `yaml:"encode"`
		TimeFormat string      `yaml:"timeFormat"`
	} `yaml:"logger"`
	ETCD struct {
		DialTimeout             time.Duration `yaml:"dialTimeout"`
		Endpoints               []string      `yaml:"endpoints"`
		Username                string        `yaml:"username"`
		Password                string        `yaml:"password"`
		HeaderTimeoutPerRequest time.Duration `yaml:"headerTimeoutPerRequest"`
		KeyPrefix               string        `yaml:"keyPrefix"`
	} `yaml:"etcd"`
}

// 加载配置文件
func loadConfigFile() error {
	var configPath string
	flag.StringVar(&configPath, "c", "./config.yaml", "配置文件路径")
	flag.Parse()
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(file).Decode(&LocalConfig)
	if err != nil {
		return err
	}
	return nil
}
