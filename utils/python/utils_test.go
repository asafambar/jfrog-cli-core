package utils

import (
	biutils "github.com/jfrog/build-info-go/utils"
	"github.com/jfrog/jfrog-cli-core/v2/xray/commands/audit/sca"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestAddRepoToPyprojectFile(t *testing.T) {
	poetryProjectPath, cleanUp := initPoetryTest(t)
	defer cleanUp()
	pyProjectPath := filepath.Join(poetryProjectPath, "pyproject.toml")
	dummyRepoName := "test-repo-name"
	dummyRepoURL := "https://ecosysjfrog.jfrog.io/"

	err := addRepoToPyprojectFile(pyProjectPath, dummyRepoName, dummyRepoURL)
	assert.NoError(t, err)
	// Validate pyproject.toml file content
	content, err := fileutils.ReadFile(pyProjectPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), dummyRepoURL)
}

func initPoetryTest(t *testing.T) (string, func()) {
	// Create and change directory to test workspace
	testAbs, err := filepath.Abs(filepath.Join("..", "..", "xray", "commands", "testdata", "poetry-project"))
	assert.NoError(t, err)
	poetryProjectPath, cleanUp := sca.CreateTestWorkspace(t, "poetry-project")
	assert.NoError(t, biutils.CopyDir(testAbs, poetryProjectPath, true, nil))
	return poetryProjectPath, cleanUp
}
