package dnf

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

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
	return run(opts, "check-update")
}

func run(opts *unipkg.Options, args ...string) error {
	command := append([]string{"dnf"}, args...)

	if opts.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	if opts.DryRun {
		command = append(command,
			"--setopt=tsflags=test",
			"--setopt=keepcache=true",
			"--cacheonly",
			"--assumeyes",
		)
	} else {
		command = append(command, "--assumeyes")
	}

	if opts.Logger != nil {
		opts.Logger(fmt.Sprintf("Running: %v", command))
	}

	cmd := exec.Command(command[0], command[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dnf failed: %w\n%s", err, string(out))
	}

	output := string(out)

	retryOut, retryErr := ensureCacheAndRetry(opts, command, output)
	if retryErr != nil {
		if opts.Logger != nil {
			opts.Logger(fmt.Sprintf("Retry after makecache failed:\n%s", string(retryOut)))
		}
		err = retryErr
	} else if retryOut != nil {
		out = retryOut
	}

	return err
}

func New() (unipkg.Manager, error) {
	if runtime.GOOS != "linux" {
		return nil, fmt.Errorf("dnf manager is only supported on Linux systems")
	}
	return &Manager{}, nil
}

func ensureCacheAndRetry(opts *unipkg.Options, command []string, output string) ([]byte, error) {
	if !opts.DryRun || !strings.Contains(strings.TrimSpace(output),
		"Cache-only enabled but no cache for repository") {
		return nil, nil
	}

	if opts.Logger != nil {
		opts.Logger("Cache missing — running 'dnf makecache' to prime metadata...")
	}

	makecache := exec.Command("sudo", "dnf", "makecache", "--setopt=keepcache=true")
	if err := makecache.Run(); err != nil {
		return nil, fmt.Errorf("failed to run makecache: %w", err)
	}

	cmd := exec.Command(command[0], command[1:]...)
	out, err := cmd.CombinedOutput()

	if opts.Logger != nil {
		opts.Logger("Retry after makecache:\n" + string(out))
	}

	return out, err
}
