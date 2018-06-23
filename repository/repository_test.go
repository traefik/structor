package repository

import (
	"log"
	"strings"
	"testing"

	"github.com/ldez/go-git-cmd-wrapper/git"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListBranches(t *testing.T) {
	git.CmdExecutor = func(name string, debug bool, args ...string) (string, error) {
		if debug {
			log.Println(name, strings.Join(args, " "))
		}
		return `
  origin/v1.3
  origin/v1.1
  origin/v1.2
`, nil
	}

	branches, err := ListBranches(true)
	require.NoError(t, err)

	expected := []string{"origin/v1.3", "origin/v1.2", "origin/v1.1"}
	assert.Equal(t, expected, branches)
}

func TestListBranches_error(t *testing.T) {
	git.CmdExecutor = func(name string, debug bool, args ...string) (string, error) {
		if debug {
			log.Println(name, strings.Join(args, " "))
		}
		return "", errors.New("fail")
	}

	_, err := ListBranches(true)
	assert.EqualError(t, err, "failed to retrieves branches: fail")
}
