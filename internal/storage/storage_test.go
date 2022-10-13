package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFullPath(t *testing.T) {
	filePath := "artifact.zip"

	fullPath := getFileFullPath(filePath)
	assert.Equal(t, fullPath, "artifact.zip")

	os.Setenv("GITHUB_SHA", "xxxxx")
	fullPath = getFileFullPath(filePath)
	assert.Equal(t, fullPath, "/github/workspace/artifact.zip")

	os.Setenv("INPUT_WORKING_DIRECTORY", "home")
	fullPath = getFileFullPath(filePath)
	assert.Equal(t, fullPath, "/github/workspace/home/artifact.zip")

	os.Setenv("INPUT_WORKING_DIRECTORY", "./home")
	fullPath = getFileFullPath(filePath)
	assert.Equal(t, fullPath, "/github/workspace/home/artifact.zip")

	os.Setenv("INPUT_WORKING_DIRECTORY", "./home/")
	fullPath = getFileFullPath(filePath)
	assert.Equal(t, fullPath, "/github/workspace/home/artifact.zip")

	os.Setenv("INPUT_WORKING_DIRECTORY", "home/")
	fullPath = getFileFullPath(filePath)
	assert.Equal(t, fullPath, "/github/workspace/home/artifact.zip")

	os.Unsetenv("GITHUB_SHA")
	fullPath = getFileFullPath(filePath)
	assert.Equal(t, fullPath, "home/artifact.zip")
}
