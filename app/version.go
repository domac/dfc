package app

import (
	"fmt"
	"runtime"
)

const Binary = "v0.0.1"

var (
	Version = fmt.Sprintf("DFS %s (build %s)", Binary, runtime.Version())
)
