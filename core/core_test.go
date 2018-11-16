package core

import (
	"os"
	"testing"

	"github.com/containous/structor/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseRequirements(t *testing.T) {
	reqts := `
pkg1<5
pkg2==3
pkg3>=1.0

pkg4-a>=1.0,<=2.0
pkg5<1.3
`

	content, err := parseRequirements([]byte(reqts))
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

func Test_getLatestReleaseTagName(t *testing.T) {
	testCases := []struct {
		desc            string
		repo            types.RepoID
		envVarLatestTag string
		expected        string
	}{
		{
			desc: "without env var override",
			repo: types.RepoID{
				Owner:          "containous",
				RepositoryName: "structor",
			},
			expected: `v\d+.\d+(.\d+)?`,
		},
		{
			desc: "with env var override",
			repo: types.RepoID{
				Owner:          "containous",
				RepositoryName: "structor",
			},
			envVarLatestTag: "foo",
			expected:        "foo",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			err := os.Setenv(envVarStructorLatestTag, test.envVarLatestTag)
			require.NoError(t, err)
			defer func() { _ = os.Unsetenv(envVarStructorLatestTag) }()

			tagName, err := getLatestReleaseTagName(test.repo)
			require.NoError(t, err)

			assert.Regexp(t, test.expected, tagName)
		})
	}
}
