package requirements

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	workingDirBasePath, err := ioutil.TempDir("", "structor-test")
	defer func() { _ = os.Remove(workingDirBasePath) }()
	require.NoError(t, err)

	testCases := []struct {
		desc                  string
		workingDirectory      string
		workingDirectoryFiles []string
		expectedErrorMessage  string
	}{
		{
			desc:                  "working case with requirements.txt in the provided directory",
			workingDirectory:      filepath.Join(workingDirBasePath, "requirements-found"),
			workingDirectoryFiles: []string{"mkdocs.yml", "requirements.txt"},
			expectedErrorMessage:  "",
		},
		{
			desc:                  "error case with no requirements.txt file found in the provided directory",
			workingDirectory:      filepath.Join(workingDirBasePath, "requirements-not-found"),
			workingDirectoryFiles: []string{"mkdocs.yml"},
			expectedErrorMessage:  "stat " + workingDirBasePath + "/requirements-not-found/requirements.txt: no such file or directory",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {

			if test.workingDirectory != "" {
				err = os.MkdirAll(test.workingDirectory, os.ModePerm)
				require.NoError(t, err)
			}

			if test.workingDirectoryFiles != nil {
				for _, repositoryFile := range test.workingDirectoryFiles {
					absoluteRepositoryFilePath := filepath.Join(test.workingDirectory, repositoryFile)
					err := os.MkdirAll(filepath.Dir(absoluteRepositoryFilePath), os.ModePerm)
					require.NoError(t, err)
					_, err = os.Create(absoluteRepositoryFilePath)
					require.NoError(t, err)
				}
			}

			resultingError := Check(test.workingDirectory)

			if test.expectedErrorMessage != "" {
				assert.EqualError(t, resultingError, test.expectedErrorMessage)
			} else {
				require.NoError(t, resultingError)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	reqts := `
pkg1<5
pkg2==3
pkg3>=1.0

pkg4-a>=1.0,<=2.0
pkg5<1.3
`

	content, err := parse([]byte(reqts))
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
