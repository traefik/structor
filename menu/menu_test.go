package menu

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/containous/structor/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildVersions(t *testing.T) {
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

func TestCreateFile(t *testing.T) {
	dir, _ := ioutil.TempDir("", "structor")
	defer os.RemoveAll(dir)

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
	buildJSFile(jsFile, versionsInfo, branches, jsTemplate)

	_, err := os.Stat(jsFile)
	require.NoError(t, err)

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
