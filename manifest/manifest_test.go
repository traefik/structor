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
				"extra": map[interface{}]interface{}{
					"social": []interface{}{
						map[interface{}]interface{}{
							"icon": "fontawesome/brands/github",
							"link": "https://github.com/drasyl-overlay/drasyl/",
							"name": "GitHub repo of drasyl",
						},
						map[interface{}]interface{}{
							"icon": "fontawesome/brands/gitlab",
							"link": "https://git.informatik.uni-hamburg.de/sane-public/drasyl",
							"name": "GitLab repo of drasyl",
						},
						map[interface{}]interface{}{
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
					map[interface{}]interface{}{"toc": map[interface{}]interface{}{"permalink": true}},
					"pymdownx.details",
					"pymdownx.inlinehilite",
					map[interface{}]interface{}{"pymdownx.highlight": map[interface{}]interface{}{"use_pygments": false}},
					"pymdownx.smartsymbols",
					map[interface{}]interface{}{"pymdownx.superfences": map[interface{}]interface{}{
						"custom_fences": []interface{}{map[interface{}]interface{}{
							"class":  "mermaid",
							"format": "!!python/name:pymdownx.superfences.fence_div_format",
							"name":   "mermaid",
						}},
					}},
					"pymdownx.tasklist",
				},
				"nav": []interface{}{
					map[interface{}]interface{}{"Welcome": "index.md"},
					map[interface{}]interface{}{
						"Getting Started": []interface{}{
							map[interface{}]interface{}{"Quick Start": "getting-started/quick-start.md"},
							map[interface{}]interface{}{"Build": "getting-started/build.md"},
							map[interface{}]interface{}{"Snapshots": "getting-started/snapshots.md"},
							map[interface{}]interface{}{"CLI": "getting-started/cli.md"},
							map[interface{}]interface{}{"Super-Peers": "getting-started/super-peers.md"},
						},
					},
					map[interface{}]interface{}{
						"Configuration": []interface{}{
							map[interface{}]interface{}{"Overview": "configuration/index.md"},
						},
					},
					map[interface{}]interface{}{"Contributing": []interface{}{
						map[interface{}]interface{}{"Submitting Issues": "contributing/submitting_issues.md"},
						map[interface{}]interface{}{"Submiting PRs": "contributing/submitting_pull_request.md"},
					}},
					map[interface{}]interface{}{
						"Architecture": []interface{}{
							map[interface{}]interface{}{"Concepts": "architecture/concepts.md"},
							map[interface{}]interface{}{"Diagrams": "architecture/diagrams.md"},
						},
					},
				},
				"plugins": []interface{}{
					"search",
					map[interface{}]interface{}{
						"git-revision-date-localized": map[interface{}]interface{}{"fallback_to_build_date": true, "type": "date"},
					},
				},
				"repo_name":        "drasyl-overlay/drasyl",
				"repo_url":         "https://github.com/drasyl-overlay/drasyl",
				"site_author":      "drasyl",
				"site_description": "drasyl Documentation",
				"site_name":        "drasyl",
				"site_url":         "https://docs.drasyl.org",
				"theme": map[interface{}]interface{}{
					"favicon": "assets/img/favicon.ico",
					"feature": map[interface{}]interface{}{"tabs": false},
					"i18n": map[interface{}]interface{}{
						"next": "Next",
						"prev": "Previous",
					},
					"icon":            map[interface{}]interface{}{"repo": "fontawesome/brands/github"},
					"include_sidebar": true,
					"language":        "en",
					"logo":            "assets/img/drasyl.png",
					"name":            "material",
					"palette": map[interface{}]interface{}{
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
					map[interface{}]interface{}{
						"toc": map[interface{}]interface{}{"permalink": true},
					},
				},
				"site_url":  "https://docs.traefik.io",
				"copyright": "Copyright &copy; 2016-2019 Containous",
				"extra": map[interface{}]interface{}{
					"traefikVersion": TempPrefixEnvName + "TRAEFIK_VERSION",
				},
				"theme": map[interface{}]interface{}{
					"include_sidebar": true,
					"favicon":         "img/traefik.icon.png",
					"feature":         map[interface{}]interface{}{"tabs": false},
					"i18n":            map[interface{}]interface{}{"prev": "Previous", "next": "Next"},
					"name":            "material",
					"custom_dir":      "docs/theme",
					"language":        "en",
					"logo":            "img/traefik.logo.png",
					"palette":         map[interface{}]interface{}{"primary": "cyan", "accent": "cyan"},
				},
				"google_analytics": []interface{}{"UA-51880359-3", "docs.traefik.io"},
				"extra_javascript": []interface{}{"theme/js/hljs/highlight.pack.js", "theme/js/extra.js"},
				"pages": []interface{}{
					map[interface{}]interface{}{
						"Getting Started": "index.md",
					},
					map[interface{}]interface{}{
						"Basics": "basics.md",
					},
					map[interface{}]interface{}{
						"Configuration": []interface{}{
							map[interface{}]interface{}{"Commons": "configuration/commons.md"},
							map[interface{}]interface{}{"Logs": "configuration/logs.md"},
							map[interface{}]interface{}{"EntryPoints": "configuration/entrypoints.md"},
							map[interface{}]interface{}{"Let's Encrypt": "configuration/acme.md"},
							map[interface{}]interface{}{"API / Dashboard": "configuration/api.md"},
							map[interface{}]interface{}{"BoltDB": "configuration/backends/boltdb.md"},
							map[interface{}]interface{}{"Consul": "configuration/backends/consul.md"},
							map[interface{}]interface{}{"Consul Catalog": "configuration/backends/consulcatalog.md"},
							map[interface{}]interface{}{"Docker": "configuration/backends/docker.md"},
							map[interface{}]interface{}{"DynamoDB": "configuration/backends/dynamodb.md"},
							map[interface{}]interface{}{"ECS": "configuration/backends/ecs.md"},
							map[interface{}]interface{}{"Etcd": "configuration/backends/etcd.md"},
							map[interface{}]interface{}{"Eureka": "configuration/backends/eureka.md"},
							map[interface{}]interface{}{"File": "configuration/backends/file.md"},
							map[interface{}]interface{}{"Kubernetes Ingress": "configuration/backends/kubernetes.md"},
							map[interface{}]interface{}{"Marathon": "configuration/backends/marathon.md"},
							map[interface{}]interface{}{"Mesos": "configuration/backends/mesos.md"},
							map[interface{}]interface{}{"Rancher": "configuration/backends/rancher.md"},
							map[interface{}]interface{}{"Rest": "configuration/backends/rest.md"},
							map[interface{}]interface{}{"Azure Service Fabric": "configuration/backends/servicefabric.md"},
							map[interface{}]interface{}{"Zookeeper": "configuration/backends/zookeeper.md"},
							map[interface{}]interface{}{"Ping": "configuration/ping.md"},
							map[interface{}]interface{}{"Metrics": "configuration/metrics.md"},
							map[interface{}]interface{}{"Tracing": "configuration/tracing.md"},
						},
					},
					map[interface{}]interface{}{
						"User Guides": []interface{}{
							map[interface{}]interface{}{"Configuration Examples": "user-guide/examples.md"},
							map[interface{}]interface{}{"Swarm Mode Cluster": "user-guide/swarm-mode.md"},
							map[interface{}]interface{}{"Swarm Cluster": "user-guide/swarm.md"},
							map[interface{}]interface{}{"Let's Encrypt & Docker": "user-guide/docker-and-lets-encrypt.md"},
							map[interface{}]interface{}{"Kubernetes": "user-guide/kubernetes.md"},
							map[interface{}]interface{}{"Marathon": "user-guide/marathon.md"},
							map[interface{}]interface{}{"Key-value Store Configuration": "user-guide/kv-config.md"},
							map[interface{}]interface{}{"Clustering/HA": "user-guide/cluster.md"},
							map[interface{}]interface{}{"gRPC Example": "user-guide/grpc.md"},
							map[interface{}]interface{}{"Traefik cluster example with Swarm": "user-guide/cluster-docker-consul.md"},
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
