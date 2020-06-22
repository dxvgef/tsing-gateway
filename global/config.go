package global

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// 引擎配置
var Config struct {
	Logger struct {
		Level      string      `yaml:"level"`
		FilePath   string      `yaml:"filePath"`
		FileMode   os.FileMode `yaml:"fileMode"`
		Encode     string      `yaml:"encode"`
		TimeFormat string      `yaml:"timeFormat"`
	} `yaml:"logger"`
	Storage struct {
		Name   string `yaml:"name"`
		Config string `yaml:"config"`
	} `yaml:"storage"`
	Proxy struct {
		IP                string        `yaml:"ip"`
		QuitWaitTimeout   time.Duration `yaml:"quitWaitTimeout"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
		HTTP              struct {
			Port uint16 `yaml:"port"`
		} `yaml:"http"`
		HTTPS struct {
			Port     uint16 `yaml:"port"`
			HTTP2    bool   `yaml:"http2"`
			CertFile string `yaml:"certFile"`
			KeyFile  string `yaml:"keyFile"`
		} `yaml:"https"`
	} `yaml:"proxy"`
	API struct {
		IP                string        `yaml:"ip"`
		Secret            string        `yaml:"secret"`
		QuitWaitTimeout   time.Duration `yaml:"quitWaitTimeout"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
		HTTP              struct {
			Port uint16 `yaml:"port"`
		} `yaml:"http"`
		HTTPS struct {
			Port     uint16 `yaml:"port"`
			HTTP2    bool   `yaml:"http2"`
			CertFile string `yaml:"certFile"`
			KeyFile  string `yaml:"keyFile"`
		} `yaml:"https"`
	} `yaml:"api"`
}

// 加载配置文件
func LoadConfigFile(configPath string) error {
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(file).Decode(&Config)
	if err != nil {
		return err
	}
	return nil
}
