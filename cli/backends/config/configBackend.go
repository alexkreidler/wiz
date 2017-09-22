package config

import (
	"github.com/tim15/wiz/api/backend"
	"github.com/tim15/wiz/cli/backends/config/json"
	"github.com/tim15/wiz/cli/backends/config/prototxt"
)

var backends []backend.ConfigBackend

func init() {
	backends = []backend.ConfigBackend{
		json.Register(),
		prototxt.Register(),
	}
}

func GetBackends() []backend.ConfigBackend {
	return backends
}
