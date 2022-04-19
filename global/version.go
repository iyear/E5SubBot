package global

import (
	"fmt"
	"runtime"
)

const (
	Version = "v1.0.0"
)

func GetRuntime() string {
	return fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
