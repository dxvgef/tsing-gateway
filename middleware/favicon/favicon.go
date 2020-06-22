package favicon

import (
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/dxvgef/tsing-gateway/global"
)

type Favicon struct {
	Status int    `json:"status"`           // HTTP状态码
	Target string `json:"target,omitempty"` // favicon.ico文件路径
}

func New(config string) (*Favicon, error) {
	var instance Favicon
	err := instance.UnmarshalJSON(global.StrToBytes(config))
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (self *Favicon) GetName() string {
	return "favicon"
}

func (self *Favicon) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	// 如果不是GET请求/favicon.ico则直接跳过
	if req.Method != "GET" {
		return true, nil
	}
	if req.RequestURI != "/favicon.ico" {
		return true, nil
	}
	// 如果是301或302请求
	if self.Status == http.StatusMovedPermanently || self.Status == http.StatusFound {
		// 使用target做为favicon.ico文件的URL
		fileURL, err := url.Parse(self.Target)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			// nolint
			resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError)))
			return false, err
		}
		resp.Header().Set("Location", fileURL.String())
		resp.WriteHeader(self.Status)
		return false, nil
	}
	if self.Status == http.StatusOK {
		fileInfo, err := os.Stat(self.Target)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			// nolint
			resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError)))
			return false, errors.New("Unable to find file '" + self.Target + "'")
		}
		if fileInfo.IsDir() {
			resp.WriteHeader(http.StatusInternalServerError)
			// nolint
			resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError)))
			return false, errors.New("`" + self.Target + "` must be a file and not a directory")
		}
		http.ServeFile(resp, req, self.Target)
		return false, nil
	}
	resp.WriteHeader(self.Status)
	return false, nil
}
