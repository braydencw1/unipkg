# 🧩 unipkg — Unified Linux/Windows Package Manager Wrapper

`unipkg` is a lightweight Go library that provides a **unified interface** for installing,
removing, and updating system packages across different operating systems.

Currently supports:

- 🐧 **APT** (Debian, Ubuntu)
- 🐧 **DNF** (Fedora, RHEL)
- 🪟 **Winget** (Windows)

---

## ✨ Features

- Common `unipkg.Manager` interface across all platforms
- Optional `sudo` usage where applicable
- Dry-run simulation mode (`Options.DryRun`)
- Built-in logging hook (`Options.Logger`)
- Smart auto-retry for DNF cache issues
- Safe cross-platform detection (`runtime.GOOS` checks)

---

## 🧱 Example Usage Auto-Detect

```go
package main

import (
	"fmt"
	"log"

	"github.com/braydencw1/unipkg"
	"github.com/braydencw1/unipkg/detect"
)

func main() {
	mgr, err := detect.GetManager()
	if err != nil {
		log.Fatalf("Failed to detect package manager: %v", err)
	}

	opts := &unipkg.Options{
		UseSudo: true,
		Logger:  func(s string) { fmt.Println("[log]", s) },
	}

	if err := mgr.Install("nginx", opts); err != nil {
		log.Fatalf("Install failed: %v", err)
	}

	fmt.Println("nginx installed successfully!")
}
```
⚙️ Example — Manual manager selection
```go
package main

import (
	"fmt"
	"log"

	"github.com/braydencw1/unipkg"
	"github.com/braydencw1/unipkg/manager/apt"
)

func main() {
	mgr := apt.New() // directly use APT backend

	opts := &unipkg.Options{
		UseSudo: true,
		DryRun:  true,
		Logger:  func(s string) { fmt.Println("[dryrun]", s) },
	}

	if err := mgr.Install("curl", opts); err != nil {
		log.Fatalf("Dry-run failed: %v", err)
	}
}
```