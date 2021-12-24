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

	"github.com/gogs/git-module"
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

func createTestRepo(t *testing.T, branch string) string {
	// figure out where we can create the test repo
	path, err := ioutil.TempDir("", "autoTagTest")
	checkFatal(t, err)

	err = os.MkdirAll(path, 0777)
	checkFatal(t, err)

	err = exec.Command("git", "init", path).Run()
	if err != nil {
		checkFatal(t, err)
	}

	// using two-step init / checkout -b to change default branch,
	// as opposed to init.defaultBranch, which would require Git 2.28+
	if branch != "" {
		err := exec.Command("git", "--git-dir="+path+"/.git", "checkout", "-b", branch).Run()
		if err != nil {
			checkFatal(t, err)
		}
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

func seedTestRepo(t *testing.T, tag string, repo *git.Repository) {
	f := repoRoot(repo) + "/README"
	err := exec.Command("touch", f).Run()
	if err != nil {
		fmt.Println("FAILED to touch the file ", f, err)
		checkFatal(t, err)
	}

	makeCommit(repo, "this is a commit")
	makeTag(repo, tag)
}

func updateReadme(t *testing.T, repo *git.Repository, content string) {
	tmpfile := repoRoot(repo) + "/README"
	err := ioutil.WriteFile(tmpfile, []byte(content), 0644)
	checkFatal(t, err)

	makeCommit(repo, content)
}

func repoRoot(r *git.Repository) string {
	checkPath := r.Path()
	if filepath.Base(checkPath) == ".git" {
		checkPath = checkPath + "/../"
	}

	p, err := filepath.Abs(checkPath)
	if err != nil {
		log.Fatal("Failed to get absolute path to repo: ", err)
	}
	return p
}
