package docker

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildImageFullName(t *testing.T) {
	testCases := []struct {
		desc             string
		imageName        string
		tagName          string
		expectedFullName string
	}{
		{
			desc:             "image and tag with no special chars",
			imageName:        "debian",
			tagName:          "slim",
			expectedFullName: "debian:slim",
		},
		{
			desc:             "image and tag with dashes",
			imageName:        "open-jdk",
			tagName:          "8-alpine",
			expectedFullName: "open-jdk:8-alpine",
		},
		{
			desc:             "image with no special chars and tag with slashes",
			imageName:        "structor",
			tagName:          "feature/antenna",
			expectedFullName: "structor:feature-antenna",
		},
		{
			desc:             "image with slashes and tag with no special chars",
			imageName:        "structor/feature",
			tagName:          "antenna",
			expectedFullName: "structor-feature:antenna",
		},
		{
			desc:             "image with column and tag with no special chars",
			imageName:        "structor:feature",
			tagName:          "antenna",
			expectedFullName: "structor-feature:antenna",
		},
		{
			desc:             "image with no special chars and tag with column",
			imageName:        "structor",
			tagName:          "feature:antenna",
			expectedFullName: "structor:feature-antenna",
		},
		{
			desc:             "Mix of everything: slashes and columns",
			imageName:        "struct:or/feat",
			tagName:          "ant/en:na",
			expectedFullName: "struct-or-feat:ant-en-na",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			dockerFullImageName := buildImageFullName(test.imageName, test.tagName)
			assert.Equal(t, dockerFullImageName, test.expectedFullName)
		})
	}
}

func TestGetDockerfile(t *testing.T) {
	workingDirBasePath, err := ioutil.TempDir("", "structor-test")
	defer func() { _ = os.Remove(workingDirBasePath) }()
	require.NoError(t, err)

	fallbackDockerfile := DockerfileInformation{
		Name:      "fallback.Dockerfile",
		ImageName: "mycompany/backend:1.2.1",
		Path:      filepath.Join(workingDirBasePath, "fallback"),
		Content:   []byte("FROM alpine:3.8"),
	}

	testCases := []struct {
		desc                 string
		workingDirectory     string
		dockerfilePath       string
		dockerfileContent    string
		dockerfileName       string
		expectedDockerfile   *DockerfileInformation
		expectedErrorMessage string
	}{
		{
			desc:              "normal case with a docs.Dockerfile at the root",
			workingDirectory:  filepath.Join(workingDirBasePath, "normal"),
			dockerfilePath:    filepath.Join(workingDirBasePath, "normal", "docs.Dockerfile"),
			dockerfileContent: "FROM alpine:3.8\n",
			dockerfileName:    "docs.Dockerfile",
			expectedDockerfile: &DockerfileInformation{
				Name:      "docs.Dockerfile",
				ImageName: "mycompany/backend:1.2.1",
				Path:      filepath.Join(workingDirBasePath, "normal", "docs.Dockerfile"),
				Content:   []byte("FROM alpine:3.8\n"),
			},
			expectedErrorMessage: "",
		},
		{
			desc:              "normal case with a docs.Dockerfile in the docs directory",
			workingDirectory:  filepath.Join(workingDirBasePath, "normal-docs"),
			dockerfilePath:    filepath.Join(workingDirBasePath, "normal-docs", "docs", "docs.Dockerfile"),
			dockerfileContent: "FROM alpine:3.8\n",
			dockerfileName:    "docs.Dockerfile",
			expectedDockerfile: &DockerfileInformation{
				Name:      "docs.Dockerfile",
				ImageName: "mycompany/backend:1.2.1",
				Path:      filepath.Join(workingDirBasePath, "normal-docs", "docs", "docs.Dockerfile"),
				Content:   []byte("FROM alpine:3.8\n"),
			},
			expectedErrorMessage: "",
		},
		{
			desc:                 "normal case with no docs.Dockerfile found",
			workingDirectory:     filepath.Join(workingDirBasePath, "normal-no-dockerfile-found"),
			dockerfilePath:       "",
			dockerfileContent:    "FROM alpine:3.8\n",
			dockerfileName:       "docs.Dockerfile",
			expectedDockerfile:   &fallbackDockerfile,
			expectedErrorMessage: "",
		},
		{
			desc:                 "error case with workingDirectory undefined",
			workingDirectory:     "",
			dockerfilePath:       filepath.Join(workingDirBasePath, "error-workingDirectory-undefined", "docs.Dockerfile"),
			dockerfileContent:    "FROM alpine:3.8\n",
			dockerfileName:       "docs.Dockerfile",
			expectedDockerfile:   nil,
			expectedErrorMessage: "workingDirectory is undefined",
		},
		{
			desc:                 "error case with workingDirectory not found",
			workingDirectory:     "not-existing",
			dockerfilePath:       filepath.Join(workingDirBasePath, "error-workingDirectory-not-found", "docs.Dockerfile"),
			dockerfileContent:    "FROM alpine:3.8\n",
			dockerfileName:       "docs.Dockerfile",
			expectedDockerfile:   nil,
			expectedErrorMessage: "stat not-existing: no such file or directory",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {

			if test.workingDirectory != "" && filepath.IsAbs(test.workingDirectory) {
				err = os.MkdirAll(test.workingDirectory, os.ModePerm)
				require.NoError(t, err)
			}

			if test.dockerfilePath != "" {
				err = os.MkdirAll(filepath.Dir(test.dockerfilePath), os.ModePerm)
				require.NoError(t, err)
				if test.dockerfileContent != "" {
					err = ioutil.WriteFile(test.dockerfilePath, []byte(test.dockerfileContent), os.ModePerm)
					require.NoError(t, err)
				}
			}

			resultingDockerfile, resultingError := GetDockerfile(test.workingDirectory, fallbackDockerfile, test.dockerfileName)

			if test.expectedErrorMessage != "" {
				assert.EqualError(t, resultingError, test.expectedErrorMessage)
			} else {
				require.NoError(t, resultingError)
				assert.Equal(t, test.expectedDockerfile, resultingDockerfile)
			}
		})
	}
}
