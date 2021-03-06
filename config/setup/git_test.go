package setup

import (
	"testing"
	"time"

	"github.com/mholt/caddy/middleware/git"
	"github.com/mholt/caddy/middleware/git/gittest"
)

// init sets the OS used to fakeOS
func init() {
	git.SetOS(gittest.FakeOS)
}

func TestGit(t *testing.T) {
	c := newTestController(`git git@github.com:mholt/caddy.git`)

	mid, err := Git(c)
	if err != nil {
		t.Errorf("Expected no errors, but got: %v", err)
	}
	if mid != nil {
		t.Fatal("Git middleware is a background service and expected to be nil.")
	}
}

func TestGitParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		expected  *git.Repo
	}{
		{`git git@github.com:user/repo`, false, &git.Repo{
			Url: "https://github.com/user/repo.git",
		}},
		{`git github.com/user/repo`, false, &git.Repo{
			Url: "https://github.com/user/repo.git",
		}},
		{`git git@github.com/user/repo`, true, nil},
		{`git http://github.com/user/repo`, false, &git.Repo{
			Url: "https://github.com/user/repo.git",
		}},
		{`git https://github.com/user/repo`, false, &git.Repo{
			Url: "https://github.com/user/repo.git",
		}},
		{`git http://github.com/user/repo {
			key ~/.key
		}`, false, &git.Repo{
			KeyPath: "~/.key",
			Url:     "git@github.com:user/repo.git",
		}},
		{`git git@github.com:user/repo {
			key ~/.key
		}`, false, &git.Repo{
			KeyPath: "~/.key",
			Url:     "git@github.com:user/repo.git",
		}},
		{`git `, true, nil},
		{`git {
		}`, true, nil},
		{`git {
		repo git@github.com:user/repo.git`, true, nil},
		{`git {
		repo git@github.com:user/repo
		key ~/.key
		}`, false, &git.Repo{
			KeyPath: "~/.key",
			Url:     "git@github.com:user/repo.git",
		}},
		{`git {
		repo git@github.com:user/repo
		key ~/.key
		interval 600
		}`, false, &git.Repo{
			KeyPath:  "~/.key",
			Url:      "git@github.com:user/repo.git",
			Interval: time.Second * 600,
		}},
		{`git {
		repo git@github.com:user/repo
		branch dev
		}`, false, &git.Repo{
			Branch: "dev",
			Url:    "https://github.com/user/repo.git",
		}},
		{`git {
		key ~/.key
		}`, true, nil},
		{`git {
		repo git@github.com:user/repo
		key ~/.key
		then echo hello world
		}`, false, &git.Repo{
			KeyPath: "~/.key",
			Url:     "git@github.com:user/repo.git",
			Then:    "echo hello world",
		}},
	}

	for i, test := range tests {
		c := newTestController(test.input)
		repo, err := gitParse(c)
		if !test.shouldErr && err != nil {
			t.Errorf("Test %v should not error but found %v", i, err)
			continue
		}
		if test.shouldErr && err == nil {
			t.Errorf("Test %v should error but found nil", i)
			continue
		}
		if !reposEqual(test.expected, repo) {
			t.Errorf("Test %v expects %v but found %v", i, test.expected, repo)
		}
	}
}

func reposEqual(expected, repo *git.Repo) bool {
	if expected == nil {
		return repo == nil
	}
	if expected.Branch != "" && expected.Branch != repo.Branch {
		return false
	}
	if expected.Host != "" && expected.Host != repo.Host {
		return false
	}
	if expected.Interval != 0 && expected.Interval != repo.Interval {
		return false
	}
	if expected.KeyPath != "" && expected.KeyPath != repo.KeyPath {
		return false
	}
	if expected.Path != "" && expected.Path != repo.Path {
		return false
	}
	if expected.Then != "" && expected.Then != repo.Then {
		return false
	}
	if expected.Url != "" && expected.Url != repo.Url {
		return false
	}
	return true
}
