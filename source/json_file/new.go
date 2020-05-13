package json_file

import (
	"encoding/json"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

type JSONFile struct {
	e          *proxy.Engine
	InputPath  string `json:"input_path"`  // 导入json文件所在的路径
	OutputPath string `json:"output_path"` // 存储导出json文件的路径
}

func New(e *proxy.Engine, config string) (*JSONFile, error) {
	var (
		err      error
		instance JSONFile
	)
	err = json.Unmarshal(global.StrToBytes(config), &instance)
	if err != nil {
		return nil, err
	}
	return &JSONFile{}, nil
}
