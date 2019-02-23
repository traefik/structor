package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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
			desc:                 "working case with mkdocs in the root of the repository",
			workingDirectory:     filepath.Join(workingDirBasePath, "mkdocs-in-root"),
			repositoryFiles:      []string{"mkdocs.yml", "requirements.txt", "docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     filepath.Join(workingDirBasePath, "mkdocs-in-root"),
			expectedErrorMessage: "",
		},
		{
			desc:                 "working case with mkdocs in ./docs",
			workingDirectory:     filepath.Join(workingDirBasePath, "mkdocs-in-docs"),
			repositoryFiles:      []string{"docs/mkdocs.yml", "docs/requirements.txt", "docs/docs.Dockerfile", ".gitignore", "docs/index.md", ".github/ISSUE.md"},
			expectedDocsRoot:     filepath.Join(workingDirBasePath, "mkdocs-in-docs", "docs"),
			expectedErrorMessage: "",
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
