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

func Test_getLatestReleaseTagName(t *testing.T) {
	testCases := []struct {
		desc                  string
		owner, repositoryName string
		envVarLatestTag       string
		expected              string
	}{
		{
			desc:           "without env var override",
			owner:          "containous",
			repositoryName: "structor",
			expected:       `v\d+.\d+(.\d+)?`,
		},
		{
			desc:            "with env var override",
			owner:           "containous",
			repositoryName:  "structor",
			envVarLatestTag: "foo",
			expected:        "foo",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			err := os.Setenv(envVarLatestTag, test.envVarLatestTag)
			require.NoError(t, err)
			defer func() { _ = os.Unsetenv(envVarLatestTag) }()

			tagName, err := getLatestReleaseTagName(test.owner, test.repositoryName)
			require.NoError(t, err)

			assert.Regexp(t, test.expected, tagName)
		})
	}
}

func Test_getDocumentationRoot(t *testing.T) {
	workingDirBasePath, err := ioutil.TempDir("", "structor-test")
	defer func() { _ = os.RemoveAll(workingDirBasePath) }()
	require.NoError(t, err)

	type expected struct {
		docsRoot string
		error    string
	}

	testCases := []struct {
		desc             string
		workingDirectory string
		repositoryFiles  []string
		expected         expected
	}{
		{
			desc:             "working case with mkdocs in the root of the repository",
			workingDirectory: filepath.Join(workingDirBasePath, "mkdocs-in-root"),
			repositoryFiles:  []string{"mkdocs.yml", "requirements.txt", "docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expected: expected{
				docsRoot: filepath.Join(workingDirBasePath, "mkdocs-in-root"),
			},
		},
		{
			desc:             "working case with mkdocs in ./docs",
			workingDirectory: filepath.Join(workingDirBasePath, "mkdocs-in-docs"),
			repositoryFiles:  []string{"docs/mkdocs.yml", "docs/requirements.txt", "docs/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expected: expected{
				docsRoot: filepath.Join(workingDirBasePath, "mkdocs-in-docs", "docs"),
			},
		},
		{
			desc:             "error case with no mkdocs file found in the search path",
			workingDirectory: filepath.Join(workingDirBasePath, "no-mkdocs-in-search-path"),
			repositoryFiles:  []string{"documentation/mkdocs.yml", "documentation/requirements.txt", "documentation/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expected: expected{
				error: "no file mkdocs.yml found in " + workingDirBasePath + "/no-mkdocs-in-search-path (search path was: /, docs/)",
			},
		},
		{
			desc:             "error case with no mkdocs file found at all",
			workingDirectory: filepath.Join(workingDirBasePath, "no-mkdocs-at-all"),
			repositoryFiles:  []string{"docs/requirements.txt", "docs/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expected: expected{
				error: "no file mkdocs.yml found in " + workingDirBasePath + "/no-mkdocs-at-all (search path was: /, docs/)",
			},
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

			docsRoot, err := getDocumentationRoot(test.workingDirectory)

			if test.expected.error != "" {
				assert.EqualError(t, err, test.expected.error)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected.docsRoot, docsRoot)
			}
		})
	}
}

func Test_getDocsDirSuffix(t *testing.T) {
	testCases := []struct {
		desc     string
		version  types.VersionsInformation
		expected string
	}{
		{
			desc: "no suffix",
			version: types.VersionsInformation{
				CurrentPath: "/tmp/structor694968263/v2.0",
				Current:     "v2.0",
			},
			expected: "",
		},
		{
			desc: "simple suffix",
			version: types.VersionsInformation{
				CurrentPath: "/tmp/structor694968263/v2.0/docs",
				Current:     "v2.0",
			},
			expected: "docs",
		},
		{
			desc: "suffix with slash",
			version: types.VersionsInformation{
				CurrentPath: "/tmp/structor694968263/v2.0/docs/",
				Current:     "v2.0",
			},
			expected: "docs",
		},
		{
			desc: "long suffix",
			version: types.VersionsInformation{
				CurrentPath: "/tmp/structor694968263/v2.0/docs/foo",
				Current:     "v2.0",
			},
			expected: "docs/foo",
		},
		{
			desc: "contains two times the version path",
			version: types.VersionsInformation{
				CurrentPath: "/tmp/structor694968263/v2.0/docs/foo/v2.0/bar",
				Current:     "v2.0",
			},
			expected: "docs/foo/v2.0/bar",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			suffix := getDocsDirSuffix(test.version)

			assert.Equal(t, test.expected, suffix)
		})
	}
}

func Test_createDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "structor-test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

	err = os.MkdirAll(filepath.Join(dir, "existing"), os.ModePerm)
	require.NoError(t, err)

	testCases := []struct {
		desc    string
		dirPath string
	}{
		{
			desc:    "when a directory doesn't already exists",
			dirPath: filepath.Join(dir, "here"),
		},
		{
			desc:    "when a directory already exists",
			dirPath: filepath.Join(dir, "existing"),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			err := createDirectory(test.dirPath)
			require.NoError(t, err)

			assert.DirExists(t, test.dirPath)
		})
	}
}
