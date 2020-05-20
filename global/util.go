package global

import (
	"errors"
	"path"
	"strings"
	"time"
	"unsafe"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

// 获得ID string值
func GetIDStr() string {
	var err error
	snowflake.Epoch = time.Now().Unix()
	IDNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Error().Msg(err.Error())
		return ""
	}
	return IDNode.Generate().String()
}

// 获得ID int64值
func GetIDInt64() int64 {
	var err error
	snowflake.Epoch = time.Now().Unix()
	IDNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Error().Msg(err.Error())
		return 0
	}
	return IDNode.Generate().Int64()
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
