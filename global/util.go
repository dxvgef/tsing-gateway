package global

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
	"sync"
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
		// 由于数据来自客户端，因此不记录日志
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
func ParseRouteFromKey(key, keyPrefix string) (hostname, routePath, routeMethod string, err error) {
	var pos int

	// 裁剪前缀
	if keyPrefix != "" {
		keyPrefix += "/routes/"
		key = strings.TrimPrefix(key, keyPrefix)
	}
	// 去掉key的第一个/，为了可以解析没有前缀的本地key
	key = strings.TrimPrefix(key, "/")

	// 查找下一个/，用于解析路由组ID
	pos = strings.Index(key, "/")
	if pos == -1 {
		err = errors.New("路由解析失败")
		// 由于数据来自客户端，因此不记录日志
		return
	}
	hostname = key[:pos]
	if hostname == "" {
		err = errors.New("主机名解析失败")
		// 由于数据来自客户端，因此不记录日志
		return
	}

	// 如果是解析存储器里的key，需要解码路由组ID
	if keyPrefix != "" {
		hostname, err = DecodeKey(hostname)
		if err != nil {
			// 由于数据来自客户端，因此不记录日志
			return
		}
	}

	// 裁剪掉key里的路由组ID部份
	key = strings.TrimPrefix(key, key[:pos+1])
	// 获取最后一次出现/符号(用于分隔路径和方法)的位置
	pos = strings.LastIndex(key, "/")
	if pos == -1 {
		err = errors.New("没有找到路径和方法的分隔符")
		// 由于数据来自客户端，因此不记录日志
		return
	}
	// 解析出路径
	routePath = key[:pos]

	// 如果是解析存储器里的key，需要解码路径
	if keyPrefix != "" {
		routePath, err = DecodeKey(routePath)
		if err != nil {
			// 由于数据来自客户端，因此不记录日志
			return
		}
	}

	// 获取方法
	routeMethod = key[pos+1:]
	if !InStr(HTTPMethods, routeMethod) {
		// 由于数据来自客户端，因此不记录日志
		err = errors.New("路由方法解析失败")
		return
	}

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

// 计算sync.Map的长度
func SyncMapLen(m *sync.Map) (count int) {
	if m == nil {
		return
	}
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return
}

// 清空sync.Map
func SyncMapClean(m *sync.Map) {
	if m == nil {
		return
	}
	m.Range(func(key, _ interface{}) bool {
		m.Delete(key)
		return true
	})
}

// 在url 添加get参数
func AddQueryValues(str string, values map[string]string) (string, error) {
	var result strings.Builder
	endpoint, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	query := endpoint.Query()
	for k := range values {
		query.Set(k, values[k])
	}
	if endpoint.Scheme != "" {
		result.WriteString(endpoint.Scheme)
		result.WriteString("://")
		result.WriteString(endpoint.Host)
		if endpoint.Path == "" {
			endpoint.Path = "/"
		}
	}
	if endpoint.Path != "" {
		result.WriteString(endpoint.Path)
	}
	result.WriteString("?")
	result.WriteString(query.Encode())
	return result.String(), nil
}
