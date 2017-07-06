package autotag

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gogits/git-module"
)

func checkFatal(t *testing.T, err error) {
	if err == nil {
		return
	}

	// The failure happens at wherever we were called, not here
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatalf("Unable to get caller")
	}
	t.Fatalf("Fail at %v:%v; %v", file, line, err)
}

func createTestRepo(t *testing.T) string {
	// figure out where we can create the test repo
	path, err := ioutil.TempDir("", "autoTagTest")
	checkFatal(t, err)

	err = os.MkdirAll(path, 0777)
	checkFatal(t, err)

	err = exec.Command("git", "init", path).Run()
	if err != nil {
		checkFatal(t, err)
	}

	tmpfile := "README"
	err = ioutil.WriteFile(path+"/"+tmpfile, []byte("foo\n"), 0644)
	checkFatal(t, err)

	return path
}

func cleanupTestRepo(t *testing.T, r *git.Repository) {
	var err error
	root := repoRoot(r)
	fmt.Println("Cleaning up test repo:", root)
	err = os.RemoveAll(root)
	checkFatal(t, err)
}

func makeCommit(r *git.Repository, msg string) {
	p := repoRoot(r)
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = p
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("make commit failed: ", string(out))
		fmt.Println(err)
	}

	cmd = exec.Command("git", "commit", "-m", msg)
	cmd.Dir = p
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("make commit failed: ", string(out))
		fmt.Println(err)
	}
}

func makeTag(r *git.Repository, tag string) {
	p := repoRoot(r)
	cmd := exec.Command("git", "tag", tag)
	cmd.Dir = p
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("tag creation failed: ", string(out))
		fmt.Println(err)
	}
}

// adds a #major comit to the repo
func newRepoMajor(t *testing.T) GitRepo {
	tr := createTestRepo(t)

	repo, err := git.OpenRepository(tr)
	checkFatal(t, err)
	seedTestRepo(t, repo)
	updateReadme(t, repo, "#major change")

	r, err := NewRepo(GitRepoConfig{RepoPath: repo.Path, Branch: "master"})
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	return *r
}

func seedTestRepo(t *testing.T, repo *git.Repository) {
	f := repoRoot(repo) + "/README"
	err := exec.Command("touch", f).Run()
	if err != nil {
		fmt.Println("FAILED to touch the file ", f, err)
		checkFatal(t, err)
	}

	makeCommit(repo, "this is a commit")
	makeTag(repo, "v1.0.1")
}

func majorTag(t *testing.T, repo *git.Repository) {
	updateReadme(t, repo, "Release version 2 #major")
	makeTag(repo, "v2.0.0")
}

func updateReadme(t *testing.T, repo *git.Repository, content string) {
	tmpfile := repoRoot(repo) + "/README"
	err := ioutil.WriteFile(tmpfile, []byte(content), 0644)
	checkFatal(t, err)

	makeCommit(repo, content)
}

func repoRoot(r *git.Repository) string {

	checkPath := r.Path
	if filepath.Base(r.Path) == ".git" {
		checkPath = r.Path + "/../"
	}

	p, err := filepath.Abs(checkPath)
	if err != nil {
		log.Fatal("Failed to get absolute path to repo: ", err)
	}
	return p
}
