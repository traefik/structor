package menu

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/structor/types"
)

func Test_buildJSFile(t *testing.T) {
	testCases := []struct {
		desc         string
		branches     []string
		versionsInfo types.VersionsInformation
		jsTemplate   string
		expected     string
	}{
		{
			desc:     "simple",
			branches: []string{"origin/v1.4", "origin/master", "v1.4.6", "origin/v1.5", "origin/v1.3"},
			versionsInfo: types.VersionsInformation{
				Current:      "v1.5",
				Latest:       "v1.4.6",
				Experimental: "master",
			},
			jsTemplate: `
var foo = [
{{- range $version := .Versions }}
	{url: "http://localhost:8080/{{ $version.Path }}", text: "{{ $version.Text }}", selected: {{ $version.Selected }} },
{{- end}}
];
`,
			expected: `
var foo = [
	{url: "http://localhost:8080/", text: "v1.4 Latest", selected: false },
	{url: "http://localhost:8080/master", text: "Experimental", selected: false },
	{url: "http://localhost:8080/v1.5", text: "v1.5 RC", selected: true },
	{url: "http://localhost:8080/v1.3", text: "v1.3", selected: false },
];
`,
		},
		{
			desc:     "sprig",
			branches: []string{"origin/v1.4", "origin/master", "v1.4.6", "origin/v1.3"},
			versionsInfo: types.VersionsInformation{
				Current:      "v1.4",
				Latest:       "v1.4.6",
				Experimental: "master",
			},
			jsTemplate: `
var foo = [
{{- range $version := .Versions }}
	{{- $text := $version.Text }}
	{{- if eq $version.State "EXPERIMENTAL" }}
		{{- $latest := semver $.Latest }}
		{{- $text = printf "v%d.%d (unreleased)" $latest.Major (int64 1 | add $latest.Minor) }}
	{{- end}}
	{url: "http://localhost:8080/{{ $version.Path }}", text: "{{ $text }}", selected: {{ eq $version.Name $.Current }} },
{{- end}}
];
`,
			expected: `
var foo = [
	{url: "http://localhost:8080/", text: "v1.4 Latest", selected: true },
	{url: "http://localhost:8080/master", text: "v1.5 (unreleased)", selected: false },
	{url: "http://localhost:8080/v1.3", text: "v1.3", selected: false },
];
`,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "structor-test")
			require.NoError(t, err)
			defer func() { _ = os.RemoveAll(dir) }()

			jsFile := filepath.Join(dir, "menu.js")

			err = buildJSFile(jsFile, test.versionsInfo, test.branches, test.jsTemplate)
			require.NoError(t, err)

			assert.FileExists(t, jsFile)

			content, err := ioutil.ReadFile(jsFile)
			require.NoError(t, err)

			assert.Equal(t, test.expected, string(content))
		})
	}
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
					Name:     "v1.4",
					State:    stateLatest,
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
					Name:     "v1.4",
					State:    stateLatest,
					Selected: true,
				},
				{
					Path:     "master",
					Text:     "Experimental",
					Name:     "master",
					State:    stateExperimental,
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
					Name:     "v1.4",
					State:    stateLatest,
					Selected: true,
				},
				{
					Path:     "v1.5",
					Text:     "v1.5 RC",
					Name:     "v1.5",
					State:    statePreFinalRelease,
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
					Name:     "v1.3",
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
					Name:     "v1.4",
					State:    stateLatest,
					Selected: true,
				},
				{
					Path:     "master",
					Text:     "Experimental",
					Name:     "master",
					State:    stateExperimental,
					Selected: false,
				},
				{
					Path:     "v1.5",
					Text:     "v1.5 RC",
					Name:     "v1.5",
					State:    statePreFinalRelease,
					Selected: false,
				},
				{
					Path:     "v1.3",
					Text:     "v1.3",
					Name:     "v1.3",
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
