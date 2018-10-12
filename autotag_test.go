package autotag

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogits/git-module"
)

func newRepo(t *testing.T, preName, preLayout string, prefix bool) GitRepo {
	path := createTestRepo(t)

	repo, err := git.OpenRepository(path)
	checkFatal(t, err)

	seedTestRepoPrefixToggle(t, repo, prefix)
	r, err := NewRepo(GitRepoConfig{
		RepoPath:                  repo.Path,
		Branch:                    "master",
		PreReleaseName:            preName,
		PreReleaseTimestampLayout: preLayout,
		Prefix:                    prefix,
	})
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	return *r
}

func newRepoWithPreReleasedTag(t *testing.T, prefix bool) GitRepo {
	path := createTestRepo(t)

	repo, err := git.OpenRepository(path)
	checkFatal(t, err)
	seedTestRepoPrefixToggle(t, repo, prefix)
	if prefix {
		makeTag(repo, "v1.0.2-pre")
	} else {
		makeTag(repo, "1.0.2-pre")
	}

	r, err := NewRepo(GitRepoConfig{RepoPath: repo.Path, Branch: "master", Prefix: prefix})
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	return *r
}

func TestBumpers(t *testing.T) {
	r := newRepo(t, "", "", true)
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
	r := newRepo(t, "", "", true)
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
	r := newRepo(t, "", "", true)
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
	expected := []string{"v1.0.2", "v1.0.1"}
	test := "TestAutoTag"

	tags := prepareRepository(t, true)
	if !compareValues(expected, tags) {
		t.Fatalf("%s expected '%+v' got '%+v'\n", test, expected, tags)
	}
}

func TestAutoTagNoPrefix(t *testing.T) {
	expected := []string{"1.0.2", "1.0.1"}
	test := "TestAutoTagNoPrefix"
	tags := prepareRepository(t, false)

	if !compareValues(expected, tags) {
		t.Fatalf("%s expected '%+v' got '%+v'\n", test, expected, tags)
	}
}

func TestAutoTagCommits(t *testing.T) {
	tags := prepareRepositoryMajor(t, true)

	expect := []string{"v2.0.0", "v1.0.1"}
	test := "TestAutoTagCommits"

	if !compareValues(expect, tags) {
		t.Fatalf("%s expected '%+v' got '%+v'\n", test, expect, tags)
	}
}

func prepareRepositoryMajor(t *testing.T, prefix bool) []string {
	r := newRepoMajorPrefixToggle(t, prefix)
	defer cleanupTestRepo(t, r.repo)
	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}
	tags, err := r.repo.GetTags()
	checkFatal(t, err)
	return tags
}

func TestAutoTagCommitsNoPrefix(t *testing.T) {
	tags := prepareRepositoryMajor(t, false)

	expect := []string{"2.0.0", "1.0.1"}
	test := "TestAutoTagCommitsNoPrefix"

	if !compareValues(expect, tags) {
		t.Fatalf("%s expected '%+v' got '%+v'\n", test, expect, tags)
	}
}

func TestAutoTagWithPreReleasedTag(t *testing.T) {
	tags := prepareRepositoryPreReleasedTag(t, true)

	expect := []string{"v1.0.2-pre", "v1.0.2", "v1.0.1"}
	test := "TestAutoTagWithPreReleasedTag"

	if !compareValues(expect, tags) {
		t.Fatalf("%s expected '%+v' got '%+v'\n", test, expect, tags)
	}
}

func TestAutoTagWithPreReleasedTagNoPrefix(t *testing.T) {
	tags := prepareRepositoryPreReleasedTag(t, false)

	test := "TestAutoTagWithPreReleasedTag"
	expect := []string{"1.0.2-pre", "1.0.2", "1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("%s expected '%+v' got '%+v'\n", test, expect, tags)
	}
}

func TestAutoTagWithPreReleaseName(t *testing.T) {
	r := newRepo(t, "test", "", true)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{"v1.0.2-test", "v1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("TestAutoTagWithPreReleaseName expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseNameNoPrefix(t *testing.T) {
	r := newRepo(t, "test", "", false)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{"1.0.2-test", "1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("TestAutoTagWithPreReleaseNameNoPrefix expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseTimestampLayout_Epoch(t *testing.T) {
	r := newRepo(t, "", "epoch", true)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("v1.0.2-%d", timeNow.Unix()), "v1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("TestAutoTagWithPreReleaseTimestampLayout_Epoch expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseTimestampLayout_EpochNoPrefix(t *testing.T) {
	r := newRepo(t, "", "epoch", false)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("1.0.2-%d", timeNow.Unix()), "1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("TestAutoTagWithPreReleaseTimestampLayout_EpochNoPrefix expected '%+v' got '%+v'\n", expect, tags)
	}
}

const testDatetimeLayout = "20060102150405"

func TestAutoTagWithPreReleaseTimestampLayout_Datetime(t *testing.T) {
	r := newRepo(t, "", testDatetimeLayout, true)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("v1.0.2-%s", timeNow.Format(testDatetimeLayout)), "v1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseTimestampLayout_DatetimeNoPrefix(t *testing.T) {
	r := newRepo(t, "", testDatetimeLayout, false)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("1.0.2-%s", timeNow.Format(testDatetimeLayout)), "1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("AutoBump expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseNameAndPreReleaseTimestampLayout(t *testing.T) {
	r := newRepo(t, "test", "epoch", true)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("v1.0.2-test.%d", timeNow.Unix()), "v1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("TestAutoTagWithPreReleaseNameAndPreReleaseTimestampLayout expected '%+v' got '%+v'\n", expect, tags)
	}
}

func TestAutoTagWithPreReleaseNameAndPreReleaseTimestampLayoutNoPrefix(t *testing.T) {
	r := newRepo(t, "test", "epoch", false)
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := []string{fmt.Sprintf("1.0.2-test.%d", timeNow.Unix()), "1.0.1"}

	if !compareValues(expect, tags) {
		t.Fatalf("TestAutoTagWithPreReleaseNameAndPreReleaseTimestampLayoutNoPrefix expected '%+v' got '%+v'\n", expect, tags)
	}
}

func compareValues(expect []string, tags []string) bool {
	found := true
	for _, val := range expect {
		found = found && hasValue(tags, val)
	}
	return found
}

func hasValue(tags []string, value string) bool {
	for _, tag := range tags {
		if tag == value {
			return true
		}
	}
	return false
}

func prepareRepository(t *testing.T, prefix bool) []string {
	r := newRepo(t, "", "", prefix)
	defer cleanupTestRepo(t, r.repo)
	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}
	tags, err := r.repo.GetTags()
	checkFatal(t, err)
	return tags
}

func prepareRepositoryPreReleasedTag(t *testing.T, prefix bool) []string {
	r := newRepoWithPreReleasedTag(t, prefix)
	defer cleanupTestRepo(t, r.repo)
	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}
	tags, err := r.repo.GetTags()
	checkFatal(t, err)
	return tags
}
