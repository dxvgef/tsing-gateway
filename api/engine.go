package api

import (
	"github.com/dxvgef/tsing-gateway/proxy"
	"github.com/dxvgef/tsing-gateway/storage"
)

var (
	sa          storage.Storage
	proxyEngine *proxy.Engine
)

func SetStorage(s storage.Storage) {
	sa = s
}

func SetProxyEngine(e *proxy.Engine) {
	proxyEngine = e
}
