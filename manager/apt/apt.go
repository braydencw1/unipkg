package apt

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/braydencw1/unipkg"
)

type Manager struct{}

var _ unipkg.Manager = (*Manager)(nil)

func (m *Manager) Install(pkg string, opts *unipkg.Options) error {
	return run(opts, "install", pkg)
}

func (m *Manager) Remove(pkg string, opts *unipkg.Options) error {
	return run(opts, "remove", pkg)
}

func (a *Manager) Update(opts *unipkg.Options) error {
	return run(opts, "upgrade")

}

func (m *Manager) Refresh(opts *unipkg.Options) error {
	return run(opts, "update")
}

func run(opts *unipkg.Options, args ...string) error {
	command := append([]string{"apt", "-y"}, args...)

	if opts.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	if opts.DryRun {
		command = append(command, "--simulate")
	}

	cmd := exec.Command(command[0], command[1:]...)

	if opts.Logger != nil {
		opts.Logger(fmt.Sprintf("Running: %v", command))
	}

	return cmd.Run()
}

func New() (unipkg.Manager, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("apt manager is only supported on Linux systems")
	}
	return &Manager{}, nil
}
