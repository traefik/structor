package requirements

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/structor/file"
	"github.com/traefik/structor/types"
)

func TestCheck(t *testing.T) {
	workingDirBasePath, err := os.MkdirTemp("", "structor-test")
	defer func() { _ = os.RemoveAll(workingDirBasePath) }()
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

func TestGetContent(t *testing.T) {
	serverURL, teardown := serveFixturesContent()
	defer teardown()

	testCases := []struct {
		desc             string
		requirementsPath string
		expected         string
	}{
		{
			desc:             "empty path",
			requirementsPath: "",
			expected:         "",
		},
		{
			desc:             "local file",
			requirementsPath: filepath.Join(".", "fixtures", "requirements.txt"),
			expected:         string(mustReadFile(filepath.Join(".", "fixtures", "requirements.txt"))),
		},
		{
			desc:             "remote file",
			requirementsPath: serverURL + "/requirements.txt",
			expected:         string(mustReadFile(filepath.Join(".", "fixtures", "requirements.txt"))),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			content, err := GetContent(test.requirementsPath)
			require.NoError(t, err)

			assert.Equal(t, test.expected, string(content))
		})
	}
}

func TestGetContent_Fail(t *testing.T) {
	serverURL, teardown := serveFixturesContent()
	defer teardown()

	testCases := []struct {
		desc             string
		requirementsPath string
	}{
		{
			desc:             "local file",
			requirementsPath: filepath.Join(".", "fixtures", "missing.txt"),
		},
		{
			desc:             "remote file",
			requirementsPath: serverURL + "/missing.txt",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			_, err := GetContent(test.requirementsPath)
			require.Error(t, err)
		})
	}
}

func TestBuild(t *testing.T) {
	testCases := []struct {
		desc          string
		customContent string
		expected      string
	}{
		{
			desc: "no custom content",
			expected: `mkdocs==0.17.5
pymdown-extensions==4.12
mkdocs-bootswatch==0.5.0
mkdocs-material==2.9.4
`,
		},
		{
			desc: "merge",
			customContent: `
mkdocs==0.17.6
`,
			expected: `mkdocs==0.17.6
mkdocs-bootswatch==0.5.0
mkdocs-material==2.9.4
pymdown-extensions==4.12
`,
		},
		{
			desc: "add",
			customContent: `
foo=0.17.6
`,
			expected: `foo=0.17.6
mkdocs==0.17.5
mkdocs-bootswatch==0.5.0
mkdocs-material==2.9.4
pymdown-extensions==4.12
`,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			dir, err := os.MkdirTemp("", "structor-test")
			require.NoError(t, err)
			defer func() { _ = os.RemoveAll(dir) }()

			requirementPath := filepath.Join(dir, filename)
			err = file.Copy(filepath.Join(".", "fixtures", filename), requirementPath)
			require.NoError(t, err)

			versionsInfo := types.VersionsInformation{
				CurrentPath: dir,
			}

			err = Build(versionsInfo, []byte(test.customContent))
			require.NoError(t, err)

			require.FileExists(t, requirementPath)
			content, err := os.ReadFile(requirementPath)
			require.NoError(t, err)

			assert.Equal(t, test.expected, string(content))
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

func mustReadFile(path string) []byte {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

func serveFixturesContent() (string, func()) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./fixtures")))

	server := httptest.NewServer(mux)

	return server.URL, server.Close
}
