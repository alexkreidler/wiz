package daemon

import (
	"github.com/tim15/wiz/api/daemon"
)

func InstallPackages(daemon.PackageList) *daemon.Status {
	return &daemon.Status{Status: true}
}
