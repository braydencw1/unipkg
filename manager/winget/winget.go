package winget

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
	return run(opts, "uninstall", pkg)
}

func (a *Manager) Update(opts *unipkg.Options) error {
	return run(opts, "upgrade", "--all")
}

func (m *Manager) Refresh(opts *unipkg.Options) error {
	return run(opts, "source", "update")
}

func run(opts *unipkg.Options, args ...string) error {
	command := append([]string{"winget", "--silent"}, args...)

	command = append(command,
		"--accept-source-agreements",
		"--accept-package-agreements",
		"--force",
	)

	// Later Implementation
	// if opts.Force {
	// 	command = append(command, "--force")
	// }

	if opts.UseSudo {
		if opts.Logger != nil {
			opts.Logger("[info] Ignoring UseSudo on Windows — not applicable")
		}
	}

	if opts.DryRun {
		if opts.Logger != nil {
			opts.Logger(fmt.Sprintf("[dry-run] Would run: %v", command))
		}
		return nil
	}

	cmd := exec.Command(command[0], command[1:]...)

	if opts.Logger != nil {
		opts.Logger(fmt.Sprintf("Running: %v", command))
	}

	return cmd.Run()
}

func New() (unipkg.Manager, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("winget manager is only supported on Windows systems")
	}
	return &Manager{}, nil
}
