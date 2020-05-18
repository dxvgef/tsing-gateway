package api

import (
	"github.com/dxvgef/tsing-gateway/storage"
)

var sa storage.Storage

func SetStorage(s storage.Storage) {
	sa = s
}
