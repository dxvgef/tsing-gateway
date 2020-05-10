package global

import (
	"os"
	"time"
)

// localConfig 全局配置
var LocalConfig struct {
	IP              string        `yaml:"ip"`
	Debug           bool          `yaml:"debug"`
	QuitWaitTimeout time.Duration `yaml:"quitWaitTimeout"`
	HTTP            struct {
		Port              int           `yaml:"port"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
	} `yaml:"http"`
	HTTPS struct {
		Port              int           `yaml:"port"`
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
	ETCD struct {
		DialTimeout          time.Duration `yaml:"dialTimeout"`
		Endpoints            []string      `yaml:"endpoints"`
		Username             string        `yaml:"username"`
		Password             string        `yaml:"password"`
		KeyPrefix            string        `yaml:"keyPrefix"`
		AutoSyncInterval     time.Duration `yaml:"autoSyncInterval"`
		DialKeepAliveTime    time.Duration `yaml:"dialKeepAliveTime"`
		DialKeepAliveTimeout time.Duration `yaml:"dialKeepAliveTimeout"`
		MaxCallSendMsgSize   int           `yaml:"maxCallSendMsgSize"`
		MaxCallRecvMsgSize   int           `yaml:"maxCallRecvMsgSize"`
		RejectOldCluster     bool          `yaml:"rejectOldCluster"`
		PermitWithoutStream  bool          `yaml:"permitWithoutStream"`
	} `yaml:"etcd"`
}
