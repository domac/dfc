package app

import (
	"fmt"
	"runtime"
)

const Binary = "v0.0.2"

var (
	Version = fmt.Sprintf("DFS %s (build %s)", Binary, runtime.Version())
)
