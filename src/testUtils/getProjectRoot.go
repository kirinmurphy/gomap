package testUtils

import (
	"path/filepath"
	"runtime"
)

func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..", "..")
}
