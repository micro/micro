package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type Gitter interface {
	Clone(repo string) error
	FetchAll(repo string) error
	Checkout(repo, branchOrCommit string) error
	RepoDir(repo string) string
}

type libGitter struct {
	folder string
}

func (g libGitter) Clone(repo string) error {
	fold := filepath.Join(g.folder, dirifyRepo(repo))
	exists, err := pathExists(fold)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = git.PlainClone(fold, false, &git.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})
	return err
}

func (g libGitter) FetchAll(repo string) error {
	repos, err := git.PlainOpen(filepath.Join(g.folder, dirifyRepo(repo)))
	if err != nil {
		return err
	}
	remotes, err := repos.Remotes()
	if err != nil {
		return err
	}

	err = remotes[0].Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func (g libGitter) Checkout(repo, branchOrCommit string) error {
	repos, err := git.PlainOpen(filepath.Join(g.folder, dirifyRepo(repo)))
	if err != nil {
		return err
	}
	worktree, err := repos.Worktree()
	if err != nil {
		return err
	}
	isCommit := func(s string) bool {
		return strings.ContainsAny(s, "0123456789") && len(s) == 40
	}
	if isCommit(branchOrCommit) {
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash:  plumbing.NewHash(branchOrCommit),
			Force: true,
		})
		if err != nil {
			return err
		}
	} else {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branchOrCommit),
			Force:  true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g libGitter) RepoDir(repo string) string {
	return filepath.Join(g.folder, dirifyRepo(repo))
}

type binaryGitter struct {
	folder string
}

func (g binaryGitter) Clone(repo string) error {
	fold := filepath.Join(g.folder, dirifyRepo(repo))
	exists, err := pathExists(fold)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	cmd := exec.Command("git", "clone", repo, ".")

	err = os.MkdirAll(fold, 0777)
	if err != nil {
		return err
	}
	cmd.Dir = fold
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	return err
}

func (g binaryGitter) FetchAll(repo string) error {
	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = filepath.Join(g.folder, dirifyRepo(repo))
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return err
}

func (g binaryGitter) Checkout(repo, branchOrCommit string) error {
	cmd := exec.Command("git", "checkout", "-f", branchOrCommit)
	cmd.Dir = filepath.Join(g.folder, dirifyRepo(repo))
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func (g binaryGitter) RepoDir(repo string) string {
	return filepath.Join(g.folder, dirifyRepo(repo))
}

func NewGitter(folder string) Gitter {
	if commandExists("git") {
		return binaryGitter{folder}
	}
	return libGitter{folder}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func dirifyRepo(s string) string {
	s = strings.ReplaceAll(s, "https://", "")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

// exists returns whether the given file or directory exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
