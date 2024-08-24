package router

import (
	"gomap/src/router"
	"gomap/src/testUtils"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	projectRoot := testUtils.GetProjectRoot()
	templateDir := filepath.Join(projectRoot, "src", "templates")

	err := os.Chdir(projectRoot)
	if err != nil {
		log.Fatal("Failed to change working directory: ", err)
	}

	router.InitializeHomePageTemplates(templateDir)
	router.InitializeUpdateMapUITemplates(templateDir)
	os.Exit(m.Run())
}
