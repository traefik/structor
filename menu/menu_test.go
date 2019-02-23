package menu

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func Test_buildVersions(t *testing.T) {
	testCases := []struct {
		desc                   string
		branches               []string
		latestTagName          string
		experimentalBranchName string
		currentVersion         string
		expected               []optionVersion
	}{
		{
			desc:           "latest",
			branches:       []string{"origin/v1.4", "v1.4.6"},
			latestTagName:  "v1.4.6",
			currentVersion: "v1.4",
			expected: []optionVersion{
				{
					Path:     "",
					Text:     "v1.4 Latest",
					Selected: true,
				},
			},
		},
		{
			desc:                   "experimental",
			branches:               []string{"origin/v1.4", "origin/master"},
			latestTagName:          "v1.4.6",
			experimentalBranchName: "master",
			currentVersion:         "v1.4",
			expected: []optionVersion{
				{
					Path:     "",
					Text:     "v1.4 Latest",
					Selected: true,
				},
				{
					Path:     "master",
					Text:     "Experimental",
					Selected: false,
				},
			},
		},
		{
			desc:           "release candidate",
			branches:       []string{"origin/v1.4", "origin/v1.5"},
			latestTagName:  "v1.4.6",
			currentVersion: "v1.4",
			expected: []optionVersion{
				{
					Path:     "",
					Text:     "v1.4 Latest",
					Selected: true,
				},
				{
					Path:     "v1.5",
					Text:     "v1.5 RC",
					Selected: false,
				},
			},
		},
		{
			desc:                   "simple version",
			branches:               []string{"origin/v1.3"},
			latestTagName:          "v1.4.6",
			experimentalBranchName: "master",
			currentVersion:         "v1.4",
			expected: []optionVersion{
				{
					Path:     "v1.3",
					Text:     "v1.3",
					Selected: false,
				},
			},
		},
		{
			desc:                   "all",
			branches:               []string{"origin/v1.4", "origin/master", "v1.4.6", "origin/v1.5", "origin/v1.3"},
			latestTagName:          "v1.4.6",
			experimentalBranchName: "master",
			currentVersion:         "v1.4",
			expected: []optionVersion{
				{
					Path:     "",
					Text:     "v1.4 Latest",
					Selected: true,
				},
				{
					Path:     "master",
					Text:     "Experimental",
					Selected: false,
				},
				{
					Path:     "v1.5",
					Text:     "v1.5 RC",
					Selected: false,
				},
				{
					Path:     "v1.3",
					Text:     "v1.3",
					Selected: false,
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			versions, err := buildVersions(test.currentVersion, test.branches, test.latestTagName, test.experimentalBranchName)
			require.NoError(t, err)

			assert.Equal(t, test.expected, versions)
		})
	}
}

func Test_buildJSFile(t *testing.T) {
	dir, _ := ioutil.TempDir("", "structor-test")
	defer func() { _ = os.RemoveAll(dir) }()

	branches := []string{"origin/v1.4", "origin/master", "v1.4.6", "origin/v1.5", "origin/v1.3"}

	versionsInfo := types.VersionsInformation{
		Current:      "v1.5",
		Latest:       "v1.4.6",
		Experimental: "master",
	}

	jsTemplate := `
var foo = [
{{- range $version := .Versions }}
	{url: "http://localhost:8080/{{ $version.Path }}", text: "{{ $version.Text }}", selected: {{ $version.Selected }} },
{{- end}}
];
`

	jsFile := filepath.Join(dir, "menu.js")

	err := buildJSFile(jsFile, versionsInfo, branches, jsTemplate)
	require.NoError(t, err)

	assert.FileExists(t, jsFile)

	content, err := ioutil.ReadFile(jsFile)
	require.NoError(t, err)

	expected := `
var foo = [
	{url: "http://localhost:8080/", text: "v1.4 Latest", selected: false },
	{url: "http://localhost:8080/master", text: "Experimental", selected: false },
	{url: "http://localhost:8080/v1.5", text: "v1.5 RC", selected: true },
	{url: "http://localhost:8080/v1.3", text: "v1.3", selected: false },
];
`

	assert.Equal(t, expected, string(content))
}

func mustReadFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}

func serveFixturesContent() (string, func()) {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./fixtures/server")))

	server := httptest.NewServer(mux)

	return server.URL, server.Close
}
