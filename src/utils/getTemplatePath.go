package utils

import "path/filepath"

func GetTemplateAbsPath(relativePath string) string {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		panic("Failed to get absolute path: " + err.Error())
	}
	return absPath
}
