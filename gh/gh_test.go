package gh

import (
	"testing"

	"github.com/containous/structor/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLatestReleaseTagName(t *testing.T) {
	repo := types.RepoID{
		Owner:          "containous",
		RepositoryName: "structor",
	}

	tagName, err := GetLatestReleaseTagName(repo)

	require.NoError(t, err)
	assert.Regexp(t, `v\d+.\d+(.\d+)?`, tagName)
}

func TestGetLatestReleaseTagName_Errors(t *testing.T) {
	repo := types.RepoID{
		Owner:          "error",
		RepositoryName: "error",
	}

	_, err := GetLatestReleaseTagName(repo)

	assert.EqualError(t, err, `failed to get latest release tag name on GitHub ("https://github.com/error/error/releases/latest"), status: 404 Not Found`)
}
