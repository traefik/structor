package repository

import (
	"sort"
	"strings"

	"github.com/ldez/go-git-cmd-wrapper/branch"
	"github.com/ldez/go-git-cmd-wrapper/git"
	gTypes "github.com/ldez/go-git-cmd-wrapper/types"
	"github.com/ldez/go-git-cmd-wrapper/worktree"
	"github.com/pkg/errors"
)

// CreateWorkTree create a worktree for a specific version
func CreateWorkTree(path string, version string, debug bool) error {
	_, err := git.Worktree(worktree.Add(path, version), git.Debugger(debug))
	if err != nil {
		return errors.Wrapf(err, "failed to add worktree on path %s for version %s", path, version)
	}

	return nil
}

// ListBranches List all remotes branches
func ListBranches(debug bool) ([]string, error) {
	branchesRaw, err := git.Branch(branch.Remotes, branch.List, branchVersionPattern, git.Debugger(debug))
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieves branches")
	}

	var branches []string
	for _, branchName := range strings.Split(branchesRaw, "\n") {
		trimmedName := strings.TrimSpace(branchName)
		if trimmedName != "" {
			branches = append(branches, trimmedName)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(branches)))
	return branches, nil
}

func branchVersionPattern(g *gTypes.Cmd) {
	g.AddOptions("origin\\/v*")
}
