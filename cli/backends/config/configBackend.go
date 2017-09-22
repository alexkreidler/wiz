package config

import (
	"github.com/tim15/wiz/api/backend"
  "github.com/tim15/wiz/api/backend/json"
)

var backends []backend.ConfigBackend

func init() {
	backends = []backend.ConfigBackend{
		{Name: "json"},
		{Name: "prototxt"},
		{Name: "hcl"}
  }
}

func GetBackends() []backend.ConfigBackend {
	return backends
}
