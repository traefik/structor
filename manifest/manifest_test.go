package manifest

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		desc     string
		filename string
		expected map[string]interface{}
	}{
		{
			desc:     "with !!python",
			filename: "sample-mkdocs.yml",
			expected: map[string]interface{}{
				"copyright": "Copyright &copy; 2020 drasyl",
				"dev_addr":  "0.0.0.0:8000",
				"docs_dir":  "content",
				"edit_uri":  "https://git.informatik.uni-hamburg.de/sane-public/drasyl/-/edit/master/docs/content",
				"extra": map[string]interface{}{
					"social": []interface{}{
						map[string]interface{}{
							"icon": "fontawesome/brands/github",
							"link": "https://github.com/drasyl-overlay/drasyl/",
							"name": "GitHub repo of drasyl",
						},
						map[string]interface{}{
							"icon": "fontawesome/brands/gitlab",
							"link": "https://git.informatik.uni-hamburg.de/sane-public/drasyl",
							"name": "GitLab repo of drasyl",
						},
						map[string]interface{}{
							"icon": "fontawesome/brands/docker",
							"link": "https://hub.docker.com/r/drasyl/drasyl",
							"name": "Docker repo of drasyl",
						},
					},
				},
				"extra_css": []interface{}{"assets/style/content.css", "assets/style/atom-one-light.css"},
				"extra_javascript": []interface{}{
					"assets/js/mermaid.min.js",
					"assets/js/hljs/highlight.min.js",
					"assets/js/extra.js",
				},
				"markdown_extensions": []interface{}{
					"admonition",
					map[string]interface{}{"toc": map[string]interface{}{"permalink": true}},
					"pymdownx.details",
					"pymdownx.inlinehilite",
					map[string]interface{}{"pymdownx.highlight": map[string]interface{}{"use_pygments": false}},
					"pymdownx.smartsymbols",
					map[string]interface{}{"pymdownx.superfences": map[string]interface{}{
						"custom_fences": []interface{}{map[string]interface{}{
							"class":  "mermaid",
							"format": "!!python/name:pymdownx.superfences.fence_div_format",
							"name":   "mermaid",
						}},
					}},
					"pymdownx.tasklist",
				},
				"nav": []interface{}{
					map[string]interface{}{"Welcome": "index.md"},
					map[string]interface{}{
						"Getting Started": []interface{}{
							map[string]interface{}{"Quick Start": "getting-started/quick-start.md"},
							map[string]interface{}{"Build": "getting-started/build.md"},
							map[string]interface{}{"Snapshots": "getting-started/snapshots.md"},
							map[string]interface{}{"CLI": "getting-started/cli.md"},
							map[string]interface{}{"Super-Peers": "getting-started/super-peers.md"},
						},
					},
					map[string]interface{}{
						"Configuration": []interface{}{
							map[string]interface{}{"Overview": "configuration/index.md"},
						},
					},
					map[string]interface{}{"Contributing": []interface{}{
						map[string]interface{}{"Submitting Issues": "contributing/submitting_issues.md"},
						map[string]interface{}{"Submiting PRs": "contributing/submitting_pull_request.md"},
					}},
					map[string]interface{}{
						"Architecture": []interface{}{
							map[string]interface{}{"Concepts": "architecture/concepts.md"},
							map[string]interface{}{"Diagrams": "architecture/diagrams.md"},
						},
					},
				},
				"plugins": []interface{}{
					"search",
					map[string]interface{}{
						"git-revision-date-localized": map[string]interface{}{"fallback_to_build_date": true, "type": "date"},
					},
				},
				"repo_name":        "drasyl-overlay/drasyl",
				"repo_url":         "https://github.com/drasyl-overlay/drasyl",
				"site_author":      "drasyl",
				"site_description": "drasyl Documentation",
				"site_name":        "drasyl",
				"site_url":         "https://docs.drasyl.org",
				"theme": map[string]interface{}{
					"favicon": "assets/img/favicon.ico",
					"feature": map[string]interface{}{"tabs": false},
					"i18n": map[string]interface{}{
						"next": "Next",
						"prev": "Previous",
					},
					"icon":            map[string]interface{}{"repo": "fontawesome/brands/github"},
					"include_sidebar": true,
					"language":        "en",
					"logo":            "assets/img/drasyl.png",
					"name":            "material",
					"palette": map[string]interface{}{
						"accent":  "teal",
						"primary": "teal",
					},
				},
			},
		},
		{
			desc:     "empty",
			filename: "empty-mkdocs.yml",
			expected: map[string]interface{}{},
		},
		{
			desc:     "traefik",
			filename: "traefik-mkdocs.yml",
			expected: map[string]interface{}{
				"site_name":        "Traefik",
				"docs_dir":         "docs",
				"dev_addr":         "0.0.0.0:8000",
				"extra_css":        []interface{}{"theme/styles/extra.css", "theme/styles/atom-one-light.css"},
				"site_description": "Traefik Documentation", "site_author": "containo.us",
				"markdown_extensions": []interface{}{
					"admonition",
					map[string]interface{}{
						"toc": map[string]interface{}{"permalink": true},
					},
				},
				"site_url":  "https://docs.traefik.io",
				"copyright": "Copyright &copy; 2016-2019 Containous",
				"extra": map[string]interface{}{
					"traefikVersion": TempPrefixEnvName + "TRAEFIK_VERSION",
				},
				"theme": map[string]interface{}{
					"include_sidebar": true,
					"favicon":         "img/traefik.icon.png",
					"feature":         map[string]interface{}{"tabs": false},
					"i18n":            map[string]interface{}{"prev": "Previous", "next": "Next"},
					"name":            "material",
					"custom_dir":      "docs/theme",
					"language":        "en",
					"logo":            "img/traefik.logo.png",
					"palette":         map[string]interface{}{"primary": "cyan", "accent": "cyan"},
				},
				"google_analytics": []interface{}{"UA-51880359-3", "docs.traefik.io"},
				"extra_javascript": []interface{}{"theme/js/hljs/highlight.pack.js", "theme/js/extra.js"},
				"pages": []interface{}{
					map[string]interface{}{
						"Getting Started": "index.md",
					},
					map[string]interface{}{
						"Basics": "basics.md",
					},
					map[string]interface{}{
						"Configuration": []interface{}{
							map[string]interface{}{"Commons": "configuration/commons.md"},
							map[string]interface{}{"Logs": "configuration/logs.md"},
							map[string]interface{}{"EntryPoints": "configuration/entrypoints.md"},
							map[string]interface{}{"Let's Encrypt": "configuration/acme.md"},
							map[string]interface{}{"API / Dashboard": "configuration/api.md"},
							map[string]interface{}{"BoltDB": "configuration/backends/boltdb.md"},
							map[string]interface{}{"Consul": "configuration/backends/consul.md"},
							map[string]interface{}{"Consul Catalog": "configuration/backends/consulcatalog.md"},
							map[string]interface{}{"Docker": "configuration/backends/docker.md"},
							map[string]interface{}{"DynamoDB": "configuration/backends/dynamodb.md"},
							map[string]interface{}{"ECS": "configuration/backends/ecs.md"},
							map[string]interface{}{"Etcd": "configuration/backends/etcd.md"},
							map[string]interface{}{"Eureka": "configuration/backends/eureka.md"},
							map[string]interface{}{"File": "configuration/backends/file.md"},
							map[string]interface{}{"Kubernetes Ingress": "configuration/backends/kubernetes.md"},
							map[string]interface{}{"Marathon": "configuration/backends/marathon.md"},
							map[string]interface{}{"Mesos": "configuration/backends/mesos.md"},
							map[string]interface{}{"Rancher": "configuration/backends/rancher.md"},
							map[string]interface{}{"Rest": "configuration/backends/rest.md"},
							map[string]interface{}{"Azure Service Fabric": "configuration/backends/servicefabric.md"},
							map[string]interface{}{"Zookeeper": "configuration/backends/zookeeper.md"},
							map[string]interface{}{"Ping": "configuration/ping.md"},
							map[string]interface{}{"Metrics": "configuration/metrics.md"},
							map[string]interface{}{"Tracing": "configuration/tracing.md"},
						},
					},
					map[string]interface{}{
						"User Guides": []interface{}{
							map[string]interface{}{"Configuration Examples": "user-guide/examples.md"},
							map[string]interface{}{"Swarm Mode Cluster": "user-guide/swarm-mode.md"},
							map[string]interface{}{"Swarm Cluster": "user-guide/swarm.md"},
							map[string]interface{}{"Let's Encrypt & Docker": "user-guide/docker-and-lets-encrypt.md"},
							map[string]interface{}{"Kubernetes": "user-guide/kubernetes.md"},
							map[string]interface{}{"Marathon": "user-guide/marathon.md"},
							map[string]interface{}{"Key-value Store Configuration": "user-guide/kv-config.md"},
							map[string]interface{}{"Clustering/HA": "user-guide/cluster.md"},
							map[string]interface{}{"gRPC Example": "user-guide/grpc.md"},
							map[string]interface{}{"Traefik cluster example with Swarm": "user-guide/cluster-docker-consul.md"},
						},
					},
				},
				"repo_name": "GitHub",
				"repo_url":  "https://github.com/traefik/traefik",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			content, err := Read(filepath.Join(".", "fixtures", test.filename))
			require.NoError(t, err)

			assert.Equal(t, test.expected, content)
		})
	}
}

func TestGetDocsDir(t *testing.T) {
	testCases := []struct {
		desc             string
		manifestFilePath string
		content          map[string]interface{}
		expected         string
	}{
		{
			desc:             "no docs_dir attribute",
			manifestFilePath: filepath.Join("foo", "bar", FileName),
			content:          map[string]interface{}{},
			expected:         filepath.Join("foo", "bar", "docs"),
		},
		{
			desc:             "with docs_dir attribute",
			manifestFilePath: filepath.Join("foo", "bar", FileName),
			content: map[string]interface{}{
				"docs_dir": "/doc",
			},
			expected: filepath.Join("foo", "bar", "doc"),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			docDir := GetDocsDir(test.content, test.manifestFilePath)

			assert.Equal(t, test.expected, docDir)
		})
	}
}

func TestAppendExtraJs(t *testing.T) {
	testCases := []struct {
		desc     string
		jsFile   string
		content  map[string]interface{}
		expected map[string]interface{}
	}{
		{
			desc:     "empty",
			jsFile:   "",
			content:  map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			desc:    "append to non existing extra_javascript attribute",
			jsFile:  "test.js",
			content: map[string]interface{}{},
			expected: map[string]interface{}{
				"extra_javascript": []interface{}{"test.js"},
			},
		},
		{
			desc:   "append to existing extra_javascript attribute",
			jsFile: "test.js",
			content: map[string]interface{}{
				"extra_javascript": []interface{}{"foo.js", "bar.js"},
			},
			expected: map[string]interface{}{
				"extra_javascript": []interface{}{"foo.js", "bar.js", "test.js"},
			},
		},
		{
			desc:   "append an already existing file to extra_javascript attribute",
			jsFile: "test.js",
			content: map[string]interface{}{
				"extra_javascript": []interface{}{"test.js"},
			},
			expected: map[string]interface{}{
				"extra_javascript": []interface{}{"test.js", "test.js"},
			},
		},
		{
			desc:   "append empty file name",
			jsFile: "",
			content: map[string]interface{}{
				"extra_javascript": []interface{}{"foo.js", "bar.js"},
			},
			expected: map[string]interface{}{
				"extra_javascript": []interface{}{"foo.js", "bar.js"},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			AppendExtraJs(test.content, test.jsFile)

			assert.Equal(t, test.expected, test.content)
		})
	}
}

func TestAppendExtraCSS(t *testing.T) {
	testCases := []struct {
		desc     string
		cssFile  string
		content  map[string]interface{}
		expected map[string]interface{}
	}{
		{
			desc:     "empty",
			cssFile:  "",
			content:  map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			desc:    "append to non existing extra_css attribute",
			cssFile: "test.css",
			content: map[string]interface{}{},
			expected: map[string]interface{}{
				"extra_css": []interface{}{"test.css"},
			},
		},
		{
			desc:    "append to existing extra_css attribute",
			cssFile: "test.css",
			content: map[string]interface{}{
				"extra_css": []interface{}{"foo.css", "bar.css"},
			},
			expected: map[string]interface{}{
				"extra_css": []interface{}{"foo.css", "bar.css", "test.css"},
			},
		},
		{
			desc:    "append an already existing file to extra_css attribute",
			cssFile: "test.css",
			content: map[string]interface{}{
				"extra_css": []interface{}{"test.css"},
			},
			expected: map[string]interface{}{
				"extra_css": []interface{}{"test.css", "test.css"},
			},
		},
		{
			desc:    "append empty file name",
			cssFile: "",
			content: map[string]interface{}{
				"extra_css": []interface{}{"foo.css", "bar.css"},
			},
			expected: map[string]interface{}{
				"extra_css": []interface{}{"foo.css", "bar.css"},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			AppendExtraCSS(test.content, test.cssFile)

			assert.Equal(t, test.expected, test.content)
		})
	}
}

func TestAddEditionURI(t *testing.T) {
	testCases := []struct {
		desc     string
		manif    map[string]interface{}
		version  string
		baseDir  string
		override bool
		expected map[string]interface{}
	}{
		{
			desc:    "no version",
			manif:   map[string]interface{}{},
			version: "",
			expected: map[string]interface{}{
				"edit_uri": "edit/master/docs/",
			},
		},
		{
			desc: "no version, no override",
			manif: map[string]interface{}{
				"edit_uri": "edit/v666/docs/",
			},
			version: "",
			expected: map[string]interface{}{
				"edit_uri": "edit/v666/docs/",
			},
		},
		{
			desc: "no version, override",
			manif: map[string]interface{}{
				"edit_uri": "edit/v1/docs/",
			},
			version:  "",
			override: true,
			expected: map[string]interface{}{
				"edit_uri": "edit/master/docs/",
			},
		},
		{
			desc: "version, no override",
			manif: map[string]interface{}{
				"edit_uri": "edit/v1/docs/",
			},
			version: "v2",
			expected: map[string]interface{}{
				"edit_uri": "edit/v1/docs/",
			},
		},
		{
			desc: "version, override",
			manif: map[string]interface{}{
				"edit_uri": "edit/v1/docs/",
			},
			version:  "v2",
			override: true,
			expected: map[string]interface{}{
				"edit_uri": "edit/v2/docs/",
			},
		},
		{
			desc: "version, no override, base dir",
			manif: map[string]interface{}{
				"edit_uri": "edit/v1/docs/",
			},
			version: "v2",
			baseDir: "foo",
			expected: map[string]interface{}{
				"edit_uri": "edit/v1/docs/",
			},
		},
		{
			desc: "version, override, base dir",
			manif: map[string]interface{}{
				"edit_uri": "edit/v1/docs/",
			},
			version:  "v2",
			baseDir:  "foo",
			override: true,
			expected: map[string]interface{}{
				"edit_uri": "edit/v2/foo/docs/",
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			AddEditionURI(test.manif, test.version, test.baseDir, test.override)

			assert.Equal(t, test.expected, test.manif)
		})
	}
}
