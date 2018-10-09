package menu

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_editManifest(t *testing.T) {
	testCases := []struct {
		desc           string
		source         string
		versionJsFile  string
		versionCSSFile string
		expected       string
	}{
		{
			desc:           "no custom files",
			source:         "fixtures/mkdocs.yml",
			versionJsFile:  "",
			versionCSSFile: "",
			expected:       "fixtures/test_no-custom-files.yml",
		},
		{
			desc:           "with JS file and existing extra",
			source:         "fixtures/mkdocs.yml",
			versionJsFile:  "structor-custom.js",
			versionCSSFile: "",
			expected:       "fixtures/test_custom-js-1.yml",
		},
		{
			desc:           "with CSS file and existing extra",
			source:         "fixtures/mkdocs.yml",
			versionJsFile:  "",
			versionCSSFile: "structor-custom.css",
			expected:       "fixtures/test_custom-css-1.yml",
		},
		{
			desc:           "with JS file and without extra",
			source:         "fixtures/mkdocs_without-extra.yml",
			versionJsFile:  "structor-custom.js",
			versionCSSFile: "",
			expected:       "fixtures/test_custom-js-2.yml",
		},
		{
			desc:           "with CSS file and without extra",
			source:         "fixtures/mkdocs_without-extra.yml",
			versionJsFile:  "",
			versionCSSFile: "structor-custom.css",
			expected:       "fixtures/test_custom-css-2.yml",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			testManifest, err := setupTestManifest(test.source)
			require.NoError(t, err)

			err = editManifest(testManifest, test.versionJsFile, test.versionCSSFile)
			require.NoError(t, err)

			assertSameContent(t, test.expected, testManifest)
		})
	}
}

func setupTestManifest(src string) (string, error) {
	srcManifest, err := filepath.Abs(src)
	if err != nil {
		return "", err
	}

	dir, err := ioutil.TempDir("", "structor")
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(dir) }()

	testManifest := filepath.Join(dir, "mkdocs.yml")

	err = fileCopy(srcManifest, testManifest)
	if err != nil {
		return "", err
	}

	return testManifest, nil
}

func assertSameContent(t *testing.T, expectedFilePath, actualFilePath string) {
	t.Helper()

	content, err := ioutil.ReadFile(actualFilePath)
	require.NoError(t, err)

	expected, err := ioutil.ReadFile(expectedFilePath)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(content))
}

func fileCopy(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = s.Close() }()

	_, err = io.Copy(f, s)
	return err
}
