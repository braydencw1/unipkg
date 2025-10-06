package detect

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/braydencw1/unipkg"
	"github.com/braydencw1/unipkg/manager/apt"
	"github.com/braydencw1/unipkg/manager/dnf"
	"github.com/braydencw1/unipkg/manager/winget"
)

var supportedPkgMgr = []string{
	"apt", "dnf", "winget",
}

func GetManager() (unipkg.Manager, error) {
	name, err := Detect()
	if err != nil {
		return nil, err
	}

	switch name {
	case "apt":
		return apt.New()
	case "dnf":
		return dnf.New()
	case "winget":
		return winget.New()
	default:
		return nil, fmt.Errorf("unsupported manager: %s", name)
	}
}

func Detect() (string, error) {
	for _, bin := range supportedPkgMgr {
		if _, err := exec.LookPath(bin); err == nil {
			return bin, nil
		}
	}
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("winget"); err == nil {
			return "winget", nil
		}
	}
	return "", fmt.Errorf("no supported package manager found")
}
