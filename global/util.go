package global

import (
	"encoding/base64"
	"errors"
	"strings"
	"sync"
	"unsafe"

	"github.com/rs/zerolog/log"
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
		log.Err(err).Str("key", key).Int("pos", pos).Caller().Send()
		return
	}
	routeGroupID = key[:pos]
	if routeGroupID == "" {
		err = errors.New("路由组ID解析失败")
		log.Err(err).Str("key", key).Int("pos", pos).Caller().Send()
		return
	}

	// 如果是解析存储器里的key，需要解码路由组ID
	if keyPrefix != "" {
		routeGroupID, err = DecodeKey(routeGroupID)
		if err != nil {
			log.Err(err).Str("routeGroupID", routeGroupID).Caller().Msg("解码路由组ID失败")
			return
		}
	}

	// 裁剪掉key里的路由组ID部份
	key = strings.TrimPrefix(key, key[:pos+1])
	// 获取最后一次出现/符号(用于分隔路径和方法)的位置
	pos = strings.LastIndex(key, "/")
	if pos == -1 {
		err = errors.New("没有找到路径和方法的分隔符")
		log.Err(err).Str("key", key).Int("pos", pos).Caller().Send()
		return
	}
	// 解析出路径
	routePath = key[:pos]

	// 如果是解析存储器里的key，需要解码路径
	if keyPrefix != "" {
		routePath, err = DecodeKey(routePath)
		if err != nil {
			log.Err(err).Str("path", routePath).Caller().Msg("路径解码失败")
			return
		}
	}

	// 获取方法
	routeMethod = key[pos+1:]
	if !InStr(HTTPMethods, routeMethod) {
		err = errors.New("路由方法解析失败")
		log.Err(err).Str("method", routeMethod).Caller().Msg("路径解码失败")
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
