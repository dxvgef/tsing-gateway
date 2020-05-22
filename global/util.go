package global

import (
	"encoding/base64"
	"errors"
	"path"
	"strings"
	"unsafe"
)

func BytesToStr(value []byte) string {
	return *(*string)(unsafe.Pointer(&value)) // nolint
}

// 编码键名
func EncodeKey(value string) string {
	return base64.RawURLEncoding.EncodeToString(StrToBytes(value))
}

// 解码键名
func DecodeKey(value string) (string, error) {
	keyBytes, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return BytesToStr(keyBytes), nil
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

// 从键名解析路由信息
func ParseRoute(key, keyPrefix string) (routeGroupID, routePath, routeMethod string, err error) {
	var pos int

	// 如果有前缀
	if keyPrefix != "" {
		keyPrefix += "/routes/"

		// 确定第一次要裁剪的位置，为了去掉前缀
		pos = strings.Index(key, keyPrefix)
		if pos == -1 {
			err = errors.New("键名前缀处理失败")
			return
		}
		key = key[pos+len(keyPrefix):]
	}

	// 去掉第一位的/字符，适用于key解析没有前缀的key
	if key[0:1] == "/" {
		key = key[1:]
	}

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

	// 解码键名
	key, err = DecodeKey(key)
	if err != nil {
		return
	}

	// 获取方法(最后一个路径)
	routeMethod = path.Base(key)
	if !InStr(Methods, routeMethod) {
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

// InStr 检查string值在一个string slice中是否存在
func InStr(s []string, str string) bool {
	for k := range s {
		if str == s[k] {
			return true
		}
	}
	return false
}
