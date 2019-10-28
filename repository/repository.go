package repository

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ldez/go-git-cmd-wrapper/branch"
	"github.com/ldez/go-git-cmd-wrapper/git"
	gTypes "github.com/ldez/go-git-cmd-wrapper/types"
	"github.com/ldez/go-git-cmd-wrapper/worktree"
)

// CreateWorkTree create a worktree for a specific version
func CreateWorkTree(path string, version string, debug bool) error {
	_, err := git.Worktree(worktree.Add(path, version), git.Debugger(debug))
	if err != nil {
		return fmt.Errorf("failed to add worktree on path %s for version %s: %w", path, version, err)
	}

	return nil
}

// ListBranches List all remotes branches
func ListBranches(debug bool) ([]string, error) {
	branchesRaw, err := git.Branch(branch.Remotes, branch.List, branchVersionPattern, git.Debugger(debug))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieves branches: %w", err)
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
