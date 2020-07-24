package global

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
)

// 参数配置
var Config struct {
	Logger struct {
		Level      string      `toml:"level"`
		FilePath   string      `toml:"filePath"`
		FileMode   os.FileMode `toml:"fileMode"`
		Encode     string      `toml:"encode"`
		TimeFormat string      `toml:"timeFormat"`
	} `toml:"logger"`
	Storage struct {
		Name   string `toml:"name"`
		Config string `toml:"config"`
	} `toml:"storage"`
	Proxy struct {
		IP                string        `toml:"ip"`
		QuitWaitTimeout   time.Duration `toml:"quitWaitTimeout"`
		ReadTimeout       time.Duration `toml:"readTimeout"`
		ReadHeaderTimeout time.Duration `toml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `toml:"writeTimeout"`
		IdleTimeout       time.Duration `toml:"idleTimeout"`
		HTTP              struct {
			Port uint16 `toml:"port"`
		} `toml:"http"`
		HTTPS struct {
			Port     uint16 `toml:"port"`
			HTTP2    bool   `toml:"http2"`
			CertFile string `toml:"certFile"`
			KeyFile  string `toml:"keyFile"`
		} `toml:"https"`
	} `toml:"proxy"`
	API struct {
		IP                string        `toml:"ip"`
		Secret            string        `toml:"secret"`
		QuitWaitTimeout   time.Duration `toml:"quitWaitTimeout"`
		ReadTimeout       time.Duration `toml:"readTimeout"`
		ReadHeaderTimeout time.Duration `toml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `toml:"writeTimeout"`
		IdleTimeout       time.Duration `toml:"idleTimeout"`
		HTTP              struct {
			Port uint16 `toml:"port"`
		} `toml:"http"`
		HTTPS struct {
			Port     uint16 `toml:"port"`
			HTTP2    bool   `toml:"http2"`
			CertFile string `toml:"certFile"`
			KeyFile  string `toml:"keyFile"`
		} `toml:"https"`
	} `toml:"api"`
}

// 加载配置文件
func LoadConfigFile(configPath string) error {
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	err = toml.NewDecoder(file).Decode(&Config)
	if err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	return nil
}
