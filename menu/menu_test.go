package menu

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/containous/structor/file"
	"github.com/containous/structor/manifest"
	"github.com/containous/structor/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTemplateContent(t *testing.T) {
	serverURL, teardown := serveFixturesContent()
	defer teardown()

	testCases := []struct {
		desc      string
		menuFiles *types.MenuFiles
		expected  Content
	}{
		{
			desc:      "no files",
			menuFiles: &types.MenuFiles{},
			expected:  Content{},
		},
		{
			desc: "JS URL",
			menuFiles: &types.MenuFiles{
				JsURL: serverURL + "/test-menu.js.gotmpl",
			},
			expected: Content{
				Js: mustReadFile("./fixtures/server/test-menu.js.gotmpl"),
			},
		},
		{
			desc: "JS URL missing file",
			menuFiles: &types.MenuFiles{
				JsURL: serverURL + "/missing-menu.js.gotmpl",
			},
			expected: Content{},
		},
		{
			desc: "JS local file",
			menuFiles: &types.MenuFiles{
				JsFile: "./fixtures/test-menu.js.gotmpl",
			},
			expected: Content{
				Js: mustReadFile("./fixtures/test-menu.js.gotmpl"),
			},
		},
		{
			desc: "JS local file missing",
			menuFiles: &types.MenuFiles{
				JsFile: "./fixtures/missing-menu.js.gotmpl",
			},
			expected: Content{},
		},
		{
			desc: "CSS URL",
			menuFiles: &types.MenuFiles{
				CSSURL: serverURL + "/test-menu.css.gotmpl",
			},
			expected: Content{
				CSS: mustReadFile("./fixtures/server/test-menu.css.gotmpl"),
			},
		},
		{
			desc: "CSS URL missing file",
			menuFiles: &types.MenuFiles{
				CSSURL: serverURL + "/missing-menu.css.gotmpl",
			},
			expected: Content{},
		},
		{
			desc: "CSS local file",
			menuFiles: &types.MenuFiles{
				CSSFile: "./fixtures/test-menu.css.gotmpl",
			},
			expected: Content{
				CSS: mustReadFile("./fixtures/test-menu.css.gotmpl"),
			},
		},
		{
			desc: "CSS local file missing",
			menuFiles: &types.MenuFiles{
				CSSFile: "./fixtures/missing-menu.css.gotmpl",
			},
			expected: Content{},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {

			content := GetTemplateContent(test.menuFiles)

			assert.Equal(t, test.expected, content)
		})
	}
}

func TestBuild(t *testing.T) {
	projectDir, err := ioutil.TempDir("", "structor-test")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(projectDir) }()

	manifestFile := filepath.Join(projectDir, manifest.FileName)
	err = file.Copy(filepath.Join(".", "fixtures", "mkdocs.yml"), manifestFile)
	require.NoError(t, err)

	versionsInfo := types.VersionsInformation{
		Latest:      "v1.7.9",
		CurrentPath: projectDir,
	}

	var branches []string

	menuContent := Content{
		Js:  mustReadFile("./fixtures/test-menu.js.gotmpl"),
		CSS: mustReadFile("./fixtures/test-menu.css.gotmpl"),
	}

	err = Build(versionsInfo, branches, menuContent)
	require.NoError(t, err)

	assert.FileExists(t, manifestFile)
	assert.FileExists(t, filepath.Join(projectDir, "docs", "theme", "js", menuJsFileName))
	assert.FileExists(t, filepath.Join(projectDir, "docs", "theme", "css", menuCSSFileName))
}
