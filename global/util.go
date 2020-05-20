package global

import (
	"errors"
	"path"
	"strings"
	"unsafe"

	"github.com/dxvgef/gommon/slice"
)

var Methods = []string{
	"*", "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT",
}

type ModuleConfig struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

func BytesToStr(value []byte) string {
	return *(*string)(unsafe.Pointer(&value)) // nolint
}

func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // nolint
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h)) // nolint
}

func FormatTime(str string) string {
	str = strings.Replace(str, "y", "2006", -1)
	str = strings.Replace(str, "m", "01", -1)
	str = strings.Replace(str, "d", "02", -1)
	str = strings.Replace(str, "h", "15", -1)
	str = strings.Replace(str, "i", "04", -1)
	str = strings.Replace(str, "s", "05", -1)
	return str
}

func ParseRoute(key, keyPrefix string) (routeGroupID, routePath, routeMethod string, err error) {
	keyPrefix += "/routes/"

	// 确定第一次要裁剪的位置，为了去掉前缀
	pos := strings.Index(key, keyPrefix)
	if pos == -1 {
		err = errors.New("键名前缀处理失败")
		return
	}
	key = key[pos+len(keyPrefix):]

	// 确定第二次要裁剪的位置，为了获得路由组ID
	pos = strings.Index(key, "/")
	if pos == -1 {
		err = errors.New("路由解析失败")
		return
	}
	routeGroupID = key[:pos]
	if routeGroupID == "" {
		err = errors.New("路由组ID解析失败")
		return
	}
	key = key[pos:]

	// 获取方法(最后一个路径)
	routeMethod = path.Base(key)
	if !slice.InStr(Methods, routeMethod) {
		err = errors.New("路由方法解析失败")
		return
	}

	// 确定第三次要裁剪的位置，为了去掉方法
	pos = strings.LastIndex(key, "/"+routeMethod)
	if pos == -1 {
		err = errors.New("获取方法的位置失败")
		return
	}
	routePath = key[:pos]
	return
}
