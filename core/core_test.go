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

func Test_findDockerFile(t *testing.T) {
	workingDirBasePath, err := ioutil.TempDir("", "structor-test")
	defer func() { _ = os.Remove(workingDirBasePath) }()
	require.NoError(t, err)

	fallbackDockerfile := dockerfileInformation{
		name:      "fallback.Dockerfile",
		imageName: "mycompany/backend:1.2.1",
		path:      filepath.Join(workingDirBasePath, "fallback"),
		content:   []byte("FROM alpine:3.8"),
	}

	testCases := []struct {
		desc                 string
		workingDirectory     string
		dockerfilePath       string
		dockerfileContent    string
		dockerfileName       string
		expectedDockerfile   *dockerfileInformation
		expectedErrorMessage string
	}{
		{
			desc:              "normal case with a docs.Dockerfile at the root",
			workingDirectory:  filepath.Join(workingDirBasePath, "normal"),
			dockerfilePath:    filepath.Join(workingDirBasePath, "normal", "docs.Dockerfile"),
			dockerfileContent: "FROM alpine:3.8\n",
			dockerfileName:    "docs.Dockerfile",
			expectedDockerfile: &dockerfileInformation{
				name:      "docs.Dockerfile",
				imageName: "mycompany/backend:1.2.1",
				path:      filepath.Join(workingDirBasePath, "normal", "docs.Dockerfile"),
				content:   []byte("FROM alpine:3.8\n"),
			},
			expectedErrorMessage: "",
		},
		{
			desc:              "normal case with a docs.Dockerfile in the docs directory",
			workingDirectory:  filepath.Join(workingDirBasePath, "normal-docs"),
			dockerfilePath:    filepath.Join(workingDirBasePath, "normal-docs", "docs", "docs.Dockerfile"),
			dockerfileContent: "FROM alpine:3.8\n",
			dockerfileName:    "docs.Dockerfile",
			expectedDockerfile: &dockerfileInformation{
				name:      "docs.Dockerfile",
				imageName: "mycompany/backend:1.2.1",
				path:      filepath.Join(workingDirBasePath, "normal-docs", "docs", "docs.Dockerfile"),
				content:   []byte("FROM alpine:3.8\n"),
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

			resultingDockerfile, resultingError := getDockerfile(fallbackDockerfile, test.workingDirectory, test.dockerfileName)

			if test.expectedErrorMessage != "" {
				assert.EqualError(t, resultingError, test.expectedErrorMessage)
			} else {
				require.NoError(t, resultingError)
				assert.Equal(t, test.expectedDockerfile, resultingDockerfile)
			}
		})
	}
}

func Test_getDocumentationRoot(t *testing.T) {
	workingDirBasePath, err := ioutil.TempDir("", "structor-test")
	defer func() { _ = os.Remove(workingDirBasePath) }()
	require.NoError(t, err)

	testCases := []struct {
		desc                 string
		workingDirectory     string
		repositoryFiles      []string
		expectedDocsRoot     string
		expectedErrorMessage string
	}{
		{
			desc:                 "working case with mkdocs and requirements.txt in the root of the repository",
			workingDirectory:     filepath.Join(workingDirBasePath, "mkdocs-req-both-in-root"),
			repositoryFiles:      []string{"mkdocs.yml", "requirements.txt", "docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     filepath.Join(workingDirBasePath, "mkdocs-req-both-in-root"),
			expectedErrorMessage: "",
		},
		{
			desc:                 "working case with mkdocs and requirements.txt both in ./docs",
			workingDirectory:     filepath.Join(workingDirBasePath, "mkdocs-req-both-in-docs"),
			repositoryFiles:      []string{"docs/mkdocs.yml", "docs/requirements.txt", "docs/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     filepath.Join(workingDirBasePath, "mkdocs-req-both-in-docs", "docs"),
			expectedErrorMessage: "",
		},
		{
			desc:                 "error case with mkdocs in the root and requirements.txt in ./docs",
			workingDirectory:     filepath.Join(workingDirBasePath, "mkdocs-in-root-req-in-docs"),
			repositoryFiles:      []string{"mkdocs.yml", "docs/requirements.txt", "docs/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     "",
			expectedErrorMessage: "stat " + workingDirBasePath + "/mkdocs-in-root-req-in-docs/requirements.txt: no such file or directory",
		},
		{
			desc:                 "error case with mkdocs in ./docs and requirements.txt in the root",
			workingDirectory:     filepath.Join(workingDirBasePath, "mkdocs-in-docs-req-in-root"),
			repositoryFiles:      []string{"docs/mkdocs.yml", "requirements.txt", "docs/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     "",
			expectedErrorMessage: "stat " + workingDirBasePath + "/mkdocs-in-docs-req-in-root/docs/requirements.txt: no such file or directory",
		},
		{
			desc:                 "error case with no mkdocs file found in the search path",
			workingDirectory:     filepath.Join(workingDirBasePath, "no-mkdocs-in-search-path"),
			repositoryFiles:      []string{"documentation/mkdocs.yml", "documentation/requirements.txt", "documentation/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     "",
			expectedErrorMessage: "no file mkdocs.yml found in " + workingDirBasePath + "/no-mkdocs-in-search-path (search path was: /,docs/)",
		},
		{
			desc:                 "error case with no mkdocs file found at all",
			workingDirectory:     filepath.Join(workingDirBasePath, "no-mkdocs-at-all"),
			repositoryFiles:      []string{"docs/requirements.txt", "docs/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     "",
			expectedErrorMessage: "no file mkdocs.yml found in " + workingDirBasePath + "/no-mkdocs-at-all (search path was: /,docs/)",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {

			if test.workingDirectory != "" {
				err = os.MkdirAll(test.workingDirectory, os.ModePerm)
				require.NoError(t, err)
			}

			if test.repositoryFiles != nil {
				for _, repositoryFile := range test.repositoryFiles {
					absoluteRepositoryFilePath := filepath.Join(test.workingDirectory, repositoryFile)
					err := os.MkdirAll(filepath.Dir(absoluteRepositoryFilePath), os.ModePerm)
					require.NoError(t, err)
					_, err = os.Create(absoluteRepositoryFilePath)
					require.NoError(t, err)
				}
			}

			resultingDocsRoot, resultingError := getDocumentationRoot(test.workingDirectory)

			if test.expectedErrorMessage != "" {
				assert.EqualError(t, resultingError, test.expectedErrorMessage)
			} else {
				require.NoError(t, resultingError)
				assert.Equal(t, test.expectedDocsRoot, resultingDocsRoot)
			}
		})
	}
}

func Test_getDockerImageFullName(t *testing.T) {

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

			dockerFullImageName := getDockerImageFullName(test.imageName, test.tagName)
			assert.Equal(t, dockerFullImageName, test.expectedFullName)
		})
	}
}
