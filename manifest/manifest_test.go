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
			desc:     "",
			filename: "empty-mkdocs.yml",
			expected: map[string]interface{}{},
		},
		{
			desc:     "",
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
				"theme": map[interface{}]interface{}{
					"include_sidebar": true,
					"favicon":         "img/traefik.icon.png",
					"feature":         map[interface{}]interface{}{"tabs": false},
					"i18n":            map[interface{}]interface{}{"prev": "Previous", "next": "Next"},
					"name":            "material",
					"custom_dir":      "docs/theme",
					"language":        "en",
					"logo":            "img/traefik.logo.png",
					"palette":         map[interface{}]interface{}{"primary": "cyan", "accent": "cyan"}},
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
				"repo_url":  "https://github.com/containous/traefik",
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

			docDir, err := GetDocsDir(test.manifestFilePath, test.content)
			require.NoError(t, err)

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

			AppendExtraJs(test.jsFile, test.content)

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

			AppendExtraCSS(test.cssFile, test.content)

			assert.Equal(t, test.expected, test.content)
		})
	}
}
