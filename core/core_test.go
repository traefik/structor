package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/containous/structor/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseRequirements(t *testing.T) {
	reqts := `
pkg1<5
pkg2==3
pkg3>=1.0

pkg4-a>=1.0,<=2.0
pkg5<1.3
`

	content, err := parseRequirements([]byte(reqts))
	require.NoError(t, err)

	expected := map[string]string{
		"pkg1":   "<5",
		"pkg2":   "==3",
		"pkg3":   ">=1.0",
		"pkg4-a": ">=1.0,<=2.0",
		"pkg5":   "<1.3",
	}

	assert.Equal(t, expected, content)
}

func Test_getLatestReleaseTagName(t *testing.T) {
	testCases := []struct {
		desc            string
		repo            types.RepoID
		envVarLatestTag string
		expected        string
	}{
		{
			desc: "without env var override",
			repo: types.RepoID{
				Owner:          "containous",
				RepositoryName: "structor",
			},
			expected: `v\d+.\d+(.\d+)?`,
		},
		{
			desc: "with env var override",
			repo: types.RepoID{
				Owner:          "containous",
				RepositoryName: "structor",
			},
			envVarLatestTag: "foo",
			expected:        "foo",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			err := os.Setenv(envVarStructorLatestTag, test.envVarLatestTag)
			require.NoError(t, err)
			defer func() { _ = os.Unsetenv(envVarStructorLatestTag) }()

			tagName, err := getLatestReleaseTagName(test.repo)
			require.NoError(t, err)

			assert.Regexp(t, test.expected, tagName)
		})
	}
}

func Test_searchAndGetDockerFile(t *testing.T) {
	workingDirPath, err := ioutil.TempDir("", "structor-test-searchAndGetDockerFile")
	defer func() { _ = os.Remove(workingDirPath) }()
	require.NoError(t, err)

	testCases := []struct {
		desc                 string
		imageName            string
		workingDirectory     string
		dockerfilePath       string
		dockerfileContent    string
		dockerfileName       string
		expectedDockerfile   dockerfileInformation
		expectedErrorMessage string
	}{
		{
			desc:              "normal case with a docs.Dockerfile at the root",
			imageName:         "mycompany/backend:1.2.1",
			workingDirectory:  filepath.Join(workingDirPath, "normal"),
			dockerfilePath:    filepath.Join(workingDirPath, "normal", "docs.Dockerfile"),
			dockerfileContent: "FROM alpine:3.8\n",
			dockerfileName:    "docs.Dockerfile",
			expectedDockerfile: dockerfileInformation{
				name:      "docs.Dockerfile",
				imageName: "mycompany/backend:1.2.1",
				path:      filepath.Join(workingDirPath, "normal", "docs.Dockerfile"),
				content:   []byte("FROM alpine:3.8\n"),
			},
			expectedErrorMessage: "",
		},
		{
			desc:              "normal case with a docs.Dockerfile in the docs directory",
			imageName:         "mycompany/backend:1.2.1",
			workingDirectory:  filepath.Join(workingDirPath, "normal-docs"),
			dockerfilePath:    filepath.Join(workingDirPath, "normal-docs", "docs", "docs.Dockerfile"),
			dockerfileContent: "FROM alpine:3.8\n",
			dockerfileName:    "docs.Dockerfile",
			expectedDockerfile: dockerfileInformation{
				name:      "docs.Dockerfile",
				imageName: "mycompany/backend:1.2.1",
				path:      filepath.Join(workingDirPath, "normal-docs", "docs", "docs.Dockerfile"),
				content:   []byte("FROM alpine:3.8\n"),
			},
			expectedErrorMessage: "",
		},
		{
			desc:                 "normal case with no docs.Dockerfile found",
			imageName:            "mycompany/backend:1.2.1",
			workingDirectory:     filepath.Join(workingDirPath, "normal-no-dockerfile-found"),
			dockerfilePath:       "",
			dockerfileContent:    "FROM alpine:3.8\n",
			dockerfileName:       "docs.Dockerfile",
			expectedDockerfile:   dockerfileInformation{},
			expectedErrorMessage: "",
		},
		{
			desc:                 "error case with no imageName provided",
			imageName:            "",
			workingDirectory:     filepath.Join(workingDirPath, "error-no-imageName"),
			dockerfilePath:       filepath.Join(workingDirPath, "error-no-imageName", "docs.Dockerfile"),
			dockerfileContent:    "FROM alpine:3.8\n",
			dockerfileName:       "docs.Dockerfile",
			expectedDockerfile:   dockerfileInformation{},
			expectedErrorMessage: "imageName is undefined",
		},
		{
			desc:                 "error case with no workingDirectory provided",
			imageName:            "mycompany/backend:1.2.1",
			workingDirectory:     "",
			dockerfilePath:       filepath.Join(workingDirPath, "error-no-workingDirectory", "docs.Dockerfile"),
			dockerfileContent:    "FROM alpine:3.8\n",
			dockerfileName:       "docs.Dockerfile",
			expectedDockerfile:   dockerfileInformation{},
			expectedErrorMessage: "workingDirectory is undefined",
		},
		{
			desc:                 "error case with workingDirectory not found",
			imageName:            "mycompany/backend:1.2.1",
			workingDirectory:     "not-existing",
			dockerfilePath:       filepath.Join(workingDirPath, "error-workingDirectory-not-found", "docs.Dockerfile"),
			dockerfileContent:    "FROM alpine:3.8\n",
			dockerfileName:       "docs.Dockerfile",
			expectedDockerfile:   dockerfileInformation{},
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
			resultingDockerfile, resultingError := findDockerfile(test.imageName, test.workingDirectory, test.dockerfileName)

			if test.expectedErrorMessage != "" {
				assert.EqualError(t, resultingError, test.expectedErrorMessage)
			} else {
				assert.NoError(t, resultingError)
			}
			assert.Equal(t, test.expectedDockerfile, resultingDockerfile)
		})
	}
}
