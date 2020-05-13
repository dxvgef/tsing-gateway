package global

import (
	"os"
	"path/filepath"
	"time"

	_ "github.com/dxvgef/filter"
	"gopkg.in/yaml.v3"
)

// 引擎配置
var Config struct {
	IP              string        `yaml:"ip"`
	Debug           bool          `yaml:"debug"`
	QuitWaitTimeout time.Duration `yaml:"quitWaitTimeout"`
	HTTP            struct {
		Port              uint          `yaml:"port"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
	} `yaml:"http"`
	HTTPS struct {
		Port              uint          `yaml:"port"`
		HTTP2             bool          `yaml:"http2"`
		CertFile          string        `yaml:"certFile"`
		KeyFile           string        `yaml:"keyFile"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
	} `yaml:"https"`
	Logger struct {
		Level      string      `yaml:"level"`
		FilePath   string      `yaml:"filePath"`
		FileMode   os.FileMode `yaml:"fileMode"`
		Encode     string      `yaml:"encode"`
		TimeFormat string      `yaml:"timeFormat"`
	} `yaml:"logger"`
	API struct {
		On     bool   `yaml:"on"`
		IP     string `yaml:"ip"`
		Port   int    `yaml:"port"`
		Path   string `yaml:"path"`
		Secret string `yaml:"secret"`
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

	if Config.HTTP.Port == 0 {
		Config.HTTP.Port = 80
	}
	if Config.HTTPS.CertFile != "" &&
		Config.HTTPS.KeyFile != "" &&
		Config.HTTPS.Port == 0 {
		Config.HTTPS.Port = 443
	}
	return nil
}
