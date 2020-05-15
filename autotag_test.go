package autotag

import (
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/gogits/git-module"
)

// testRepoSetup provides a method for setting up a new git repo in a temporary directory
type testRepoSetup struct {
	// (optional) versioning scheme to use, eg: "" or "autotag", "conventional". If not set, defaults to "" (autotag)
	scheme string

	// (optional) branch to create. If not set, defaults to "master"
	branch string

	// (optional) initial tag. If not set, defaults to "v0.0.1"
	initialTag string

	// (optional) extra tags to add to the repo
	extraTags []string

	// (optional) the prerelease name to use, eg "pre". If not set, no prerelease name will be used
	preReleaseName string

	// (optional) the prerelease timestamp format to use, eg: "epoch". If not set, no prerelease timestamp will be used
	preReleaseTimestampLayout string

	// (optional) commit message to use for the next, untagged commit. Settings this allows for testing the
	// commit message parsing logic. eg: "#major this is a major commit"
	nextCommit string
}

// newTestRepo creates a new git repo in a temporary directory and returns an autotag.GitRepo struct for
// testing the autotag package.
// You must call cleanupTestRepo(t, r.repo) to remove the temporary directory after running tests.
func newTestRepo(t *testing.T, setup testRepoSetup) GitRepo {
	tr := createTestRepo(t)

	repo, err := git.OpenRepository(tr)
	checkFatal(t, err)

	branch := setup.branch
	if branch == "" {
		branch = "master"
	}

	tag := setup.initialTag
	if setup.initialTag == "" {
		tag = "v0.0.1"
	}
	seedTestRepo(t, tag, repo)

	if len(setup.extraTags) > 0 {
		for _, t := range setup.extraTags {
			makeTag(repo, t)
		}
	}

	if setup.nextCommit != "" {
		updateReadme(t, repo, setup.nextCommit)
	}

	r, err := NewRepo(GitRepoConfig{
		RepoPath:                  repo.Path,
		Branch:                    branch,
		PreReleaseName:            setup.preReleaseName,
		PreReleaseTimestampLayout: setup.preReleaseTimestampLayout,
		Scheme:                    setup.scheme,
	})
	if err != nil {
		t.Fatal("Error creating repo", err)
	}

	return *r
}

func TestMajor(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag: "v1.0.1",
	})
	defer cleanupTestRepo(t, r.repo)

	v, err := r.MajorBump()
	if err != nil {
		t.Fatal("MajorBump failed: ", err)
	}

	if v.String() != "2.0.0" {
		t.Fatalf("MajorBump failed expected '2.0.0' got '%s' ", v)
	}

	fmt.Printf("Major is now %s\n", v)
}

func TestMinor(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag: "v1.0.1",
	})
	defer cleanupTestRepo(t, r.repo)

	v, err := r.MinorBump()
	if err != nil {
		t.Fatal("MinorBump failed: ", err)
	}

	if v.String() != "1.1.0" {
		t.Fatalf("MinorBump failed expected '1.1.0' got '%s' \n", v)
	}
}

func TestPatch(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag: "v1.0.1",
	})
	defer cleanupTestRepo(t, r.repo)

	v, err := r.PatchBump()
	if err != nil {
		t.Fatal("PatchBump failed: ", err)
	}

	if v.String() != "1.0.2" {
		t.Fatalf("PatchBump failed expected '1.0.2' got '%s' \n", v)
	}
}

func TestAutoTag(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag: "v1.0.1",
	})
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)
	assert.Contains(t, tags, "v1.0.2")
}

func TestAutoTagCommits(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag: "v1.0.1",
		nextCommit: "#major change",
	})
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)
	assert.Contains(t, tags, "v2.0.0")
}

func TestAutoTagWithPreReleasedTag(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag: "v1.0.1",
		extraTags:  []string{"v1.0.2-pre"},
	})
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	assert.Contains(t, tags, "v1.0.2-pre")
	assert.Contains(t, tags, "v1.0.2")
}

func TestAutoTagWithPreReleaseName(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag:     "v1.0.1",
		preReleaseName: "test",
	})
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	assert.Contains(t, tags, "v1.0.2-test")
}

func TestAutoTagWithPreReleaseTimestampLayout_Epoch(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag:                "v1.0.1",
		preReleaseTimestampLayout: "epoch",
	})
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := fmt.Sprintf("v1.0.2-%d", timeNow.Unix())
	assert.Contains(t, tags, expect)
}

const testDatetimeLayout = "20060102150405"

func TestAutoTagWithPreReleaseTimestampLayout_Datetime(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag:                "v1.0.1",
		preReleaseTimestampLayout: testDatetimeLayout,
	})
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := fmt.Sprintf("v1.0.2-%s", timeNow.Format(testDatetimeLayout))
	assert.Contains(t, tags, expect)
}

func TestAutoTagWithPreReleaseNameAndPreReleaseTimestampLayout(t *testing.T) {
	r := newTestRepo(t, testRepoSetup{
		initialTag:                "v1.0.1",
		preReleaseName:            "test",
		preReleaseTimestampLayout: "epoch",
	})
	defer cleanupTestRepo(t, r.repo)

	err := r.AutoTag()
	timeNow := time.Now().UTC()

	if err != nil {
		t.Fatal("AutoTag failed ", err)
	}

	tags, err := r.repo.GetTags()
	checkFatal(t, err)

	expect := fmt.Sprintf("v1.0.2-test.%d", timeNow.Unix())
	assert.Contains(t, tags, expect)
}

func TestParseCommit(t *testing.T) {
	tests := []struct {
		name        string
		scheme      string
		nextCommit  string
		initialTag  string
		expectedTag string
	}{
		// tests for autotag (default) scheme
		{
			name:        "autotag scheme, [major] bump",
			scheme:      "autotag",
			nextCommit:  "[major] this is a big release\n\nfoo bar baz\n",
			initialTag:  "v1.0.0",
			expectedTag: "v2.0.0",
		},
		{
			name:        "autotag scheme, [minor] bump",
			scheme:      "autotag",
			nextCommit:  "[minor] this is a smaller release\n\nfoo bar baz\n",
			initialTag:  "v1.0.0",
			expectedTag: "v1.1.0",
		},
		{
			name:        "autotag scheme, patch bump",
			scheme:      "autotag",
			nextCommit:  "this is just a basic change\n\nfoo bar baz\n",
			initialTag:  "v1.0.0",
			expectedTag: "v1.0.1",
		},
		{
			name:        "autotag scheme, #major bump",
			scheme:      "autotag",
			nextCommit:  "#major this is a big release\n\nfoo bar baz\n",
			initialTag:  "v1.0.0",
			expectedTag: "v2.0.0",
		},
		{
			name:        "autotag scheme, #minor bump",
			scheme:      "autotag",
			nextCommit:  "#minor this is a smaller release\n\nfoo bar baz\n",
			initialTag:  "v1.0.0",
			expectedTag: "v1.1.0",
		},
		// tests for conventional commits scheme. Based on:
		// https://www.conventionalcommits.org/en/v1.0.0/#summary
		// and
		// https://www.conventionalcommits.org/en/v1.0.0/#examples
		{
			name:        "conventional commits, minor bump without scope",
			scheme:      "conventional",
			nextCommit:  "feat: allow provided config object to extend other configs",
			initialTag:  "v1.0.0",
			expectedTag: "v1.1.0",
		},
		{
			name:        "conventional commits, minor bump with scope",
			scheme:      "conventional",
			nextCommit:  "feat(lang): add polish language",
			initialTag:  "v1.0.0",
			expectedTag: "v1.1.0",
		},
		{
			name:        "conventional commits, breaking change via ! appended to type",
			scheme:      "conventional",
			nextCommit:  "refactor!: drop support for Node 6",
			initialTag:  "v1.0.0",
			expectedTag: "v2.0.0",
		},
		{
			name:        "conventional commits, breaking change via ! appended to type/scope",
			scheme:      "conventional",
			nextCommit:  "refactor(runtime)!: drop support for Node 6",
			initialTag:  "v1.0.0",
			expectedTag: "v2.0.0",
		},
		{
			name:        "conventional commits, breaking change via footer",
			scheme:      "conventional",
			nextCommit:  "feat: allow provided config object to extend other configs\n\nbody before footer\n\nBREAKING CHANGE: non-backwards compatible",
			initialTag:  "v1.0.0",
			expectedTag: "v2.0.0",
		},
		{
			name:        "conventional commits, patch/minor bump",
			scheme:      "conventional",
			nextCommit:  "fix: correct minor typos in code",
			initialTag:  "v1.0.0",
			expectedTag: "v1.0.1",
		},
		{
			name:        "conventional commits, non-conforming fallback to patch bump",
			scheme:      "conventional",
			nextCommit:  "not a conventional commit message",
			initialTag:  "v1.0.0",
			expectedTag: "v1.0.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestRepo(t, testRepoSetup{
				scheme:     tc.scheme,
				initialTag: tc.initialTag,
				nextCommit: tc.nextCommit,
			})
			defer cleanupTestRepo(t, r.repo)

			err := r.AutoTag()
			if err != nil {
				t.Fatal("AutoTag failed ", err)
			}

			tags, err := r.repo.GetTags()
			checkFatal(t, err)
			assert.Contains(t, tags, tc.expectedTag)
		})
	}
}
