package autotag

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gogits/git-module"
)

func newRepo(t *testing.T, preName, preLayout string) GitRepo {
	path := createTestRepo(t)

	repo, err := git.OpenRepository(path)
	checkFatal(t, err)

	seedTestRepo(t, repo)
	r, err := NewRepo(GitRepoConfig{
		RepoPath:                  repo.Path,
		Branch:                    "master",
		PreReleaseName:            preName,
		PreReleaseTimestampLayout: preLayout,
	})
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	return *r
}

func newRepoWithPreReleasedTag(t *testing.T) GitRepo {
	path := createTestRepo(t)

	repo, err := git.OpenRepository(path)
	checkFatal(t, err)
	seedTestRepo(t, repo)
	makeTag(repo, "v1.0.2-pre")

	r, err := NewRepo(GitRepoConfig{RepoPath: repo.Path, Branch: "master"})
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	return *r
}

func TestBumpers(t *testing.T) {
	r := newRepo(t, "", "")
	defer cleanupTestRepo(t, r.repo)

	majorTag(t, r.repo)
	v, err := r.MajorBump()
	if err != nil {
		t.Fatal("MajorBump failed: ", err)
	}

	if v.String() != "2.0.0" {
		t.Fatalf("MajorBump failed expected '2.0.1' got '%s' ", v)
	}

	fmt.Printf("Major is now %s\n", v)
}
func TestMinor(t *testing.T) {
	r := newRepo(t, "", "")
	defer cleanupTestRepo(t, r.repo)

	majorTag(t, r.repo)
	v, err := r.MinorBump()
	if err != nil {
		t.Fatal("MinorBump failed: ", err)
	}

	if v.String() != "1.1.0" {
		t.Fatalf("MinorBump failed expected '1.1.0' got '%s' \n", v)
	}
}
func TestPatch(t *testing.T) {
	r := newRepo(t, "", "")
	defer cleanupTestRepo(t, r.repo)

	majorTag(t, r.repo)
	v, err := r.PatchBump()
	if err != nil {
		t.Fatal("PatchBump failed: ", err)
	}

	if v.String() != "1.0.2" {
		t.Fatalf("PatchBump failed expected '1.0.2' got '%s' \n", v)
	}
}

func TestAutoTag(t *testing.T) {
	r := newRepo(t, "", "")
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	// XXX(fujin): When switching to `git-module`, the sort order reversed. Most recent is first.
	expect := []string{"v1.0.2", "v1.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagCommits(t *testing.T) {
	r := newRepoMajor(t)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)
	// XXX(fujin): When switching to `git-module`, the sort order reversed. Most recent is first.
	expect := []string{"v2.0.0", "v1.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleasedTag(t *testing.T) {
	r := newRepoWithPreReleasedTag(t)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{"v1.0.2-pre", "v1.0.2", "v1.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseName(t *testing.T) {
	r := newRepo(t, "test", "")
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{"v1.0.2-test", "v1.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseTimestampLayout_Epoch(t *testing.T) {
	r := newRepo(t, "", "epoch")
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("v1.0.2-%d", timeNow.Unix()), "v1.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

const testDatetimeLayout = "20060102150405"

func TestAutoTagWithPreReleaseTimestampLayout_Datetime(t *testing.T) {
	r := newRepo(t, "", testDatetimeLayout)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("v1.0.2-%s", timeNow.Format(testDatetimeLayout)), "v1.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseNameAndPreReleaseTimestampLayout(t *testing.T) {
	r := newRepo(t, "test", "epoch")
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("v1.0.2-test.%d", timeNow.Unix()), "v1.0.1"}

	if !reflect.DeepEqual(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}
