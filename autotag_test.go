package autotag

import (
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/gogits/git-module"
)

func init() {
	// fixed point-in-time time.Now() for testing
	timeNow = func() time.Time {
		return time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	}
}

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

	// (optional) build metadata to append to the version
	buildMetadata string

	// (optional) prepend literal 'v' to version tags (default: true)
	disablePrefix bool

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
		if setup.disablePrefix {
			tag = "0.0.1"
		}
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
		BuildMetadata:             setup.buildMetadata,
		Scheme:                    setup.scheme,
		Prefix:                    !setup.disablePrefix,
	})

	if err != nil {
		t.Fatal("Error creating repo: ", err)
	}

	return *r
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       GitRepoConfig
		shouldErr bool
	}{
		{
			name:      "missing branch",
			cfg:       GitRepoConfig{},
			shouldErr: true,
		},
		{
			name: "invalid build metadata",
			cfg: GitRepoConfig{
				Branch:        "master",
				BuildMetadata: "foo..bar",
			},
			shouldErr: true,
		},
		{
			name: "invalid build metadata - purely empty identifier",
			cfg: GitRepoConfig{
				Branch:        "master",
				BuildMetadata: "...",
			},
			shouldErr: true,
		},
		{
			name: "invalid pre-release-name - leading zero",
			cfg: GitRepoConfig{
				Branch:         "master",
				PreReleaseName: "024",
			},
			shouldErr: true,
		},
		{
			name: "invalid pre-release-name - empty identifier",
			cfg: GitRepoConfig{
				Branch:         "master",
				PreReleaseName: "...",
			},
			shouldErr: true,
		},
		{
			name: "invalid pre-release-timestamp",
			cfg: GitRepoConfig{
				Branch:                    "master",
				PreReleaseTimestampLayout: "foo",
			},
			shouldErr: true,
		},
		{
			name: "valid config with all options used",
			cfg: GitRepoConfig{
				Branch:                    "master",
				PreReleaseName:            "foo",
				PreReleaseTimestampLayout: "epoch",
				BuildMetadata:             "g12345678",
				Prefix:                    true,
			},
			shouldErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateConfig(tc.cfg)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
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

func TestMissingInitialTag(t *testing.T) {
	tr := createTestRepo(t)
	repo, err := git.OpenRepository(tr)
	checkFatal(t, err)
	defer cleanupTestRepo(t, repo)

	updateReadme(t, repo, "a commit before any usable tag has been created")

	_, err = NewRepo(GitRepoConfig{
		RepoPath: repo.Path,
		Branch:   "master",
	})
	assert.Error(t, err)
}

func TestAutoTag(t *testing.T) {
	tests := []struct {
		name        string
		setup       testRepoSetup
		shouldErr   bool
		expectedTag string
	}{
		// tests for autotag (default) scheme
		{
			name: "autotag scheme, [major] bump",
			setup: testRepoSetup{
				scheme:     "autotag",
				nextCommit: "[major] this is a big release\n\nfoo bar baz\n",
				initialTag: "v1.0.0",
			},
			expectedTag: "v2.0.0",
		},
		{
			name: "autotag scheme, [minor] bump",
			setup: testRepoSetup{
				scheme:     "autotag",
				nextCommit: "[minor] this is a smaller release\n\nfoo bar baz\n",
				initialTag: "v1.0.0",
			},
			expectedTag: "v1.1.0",
		},
		{
			name: "autotag scheme, patch bump",
			setup: testRepoSetup{
				scheme:     "autotag",
				nextCommit: "this is just a basic change\n\nfoo bar baz\n",
				initialTag: "v1.0.0",
			},
			expectedTag: "v1.0.1",
		},
		{
			name: "autotag scheme, #major bump",
			setup: testRepoSetup{
				scheme:     "autotag",
				nextCommit: "#major this is a big release\n\nfoo bar baz\n",
				initialTag: "v1.0.0",
			},
			expectedTag: "v2.0.0",
		},
		{
			name: "autotag scheme, #minor bump",
			setup: testRepoSetup{
				scheme:     "autotag",
				nextCommit: "#minor this is a smaller release\n\nfoo bar baz\n",
				initialTag: "v1.0.0",
			},
			expectedTag: "v1.1.0",
		},
		{
			name: "pre-release-name with patch bump",
			setup: testRepoSetup{
				scheme:         "autotag",
				nextCommit:     "#patch bump",
				initialTag:     "v1.0.0",
				preReleaseName: "dev",
			},
			expectedTag: "v1.0.1-dev",
		},
		{
			name: "epoch pre-release-timestamp",
			setup: testRepoSetup{
				scheme:                    "autotag",
				nextCommit:                "#patch bump",
				initialTag:                "v1.0.0",
				preReleaseTimestampLayout: "epoch",
			},
			expectedTag: fmt.Sprintf("v1.0.1-%d", timeNow().UTC().Unix()),
		},
		{
			name: "datetime pre-release-timestamp",
			setup: testRepoSetup{
				scheme:                    "autotag",
				nextCommit:                "#patch bump",
				initialTag:                "v1.0.0",
				preReleaseTimestampLayout: "datetime",
			},
			expectedTag: fmt.Sprintf("v1.0.1-%s", timeNow().Format(datetimeTsLayout)),
		},
		{
			name: "epoch pre-release-timestamp and pre-release-name",
			setup: testRepoSetup{
				scheme:                    "autotag",
				nextCommit:                "#patch bump",
				initialTag:                "v1.0.0",
				preReleaseName:            "dev",
				preReleaseTimestampLayout: "epoch",
			},
			expectedTag: fmt.Sprintf("v1.0.1-dev.%d", timeNow().UTC().Unix()),
		},
		{
			name: "ignore non-conforming tags",
			setup: testRepoSetup{
				scheme:     "autotag",
				nextCommit: "#patch bump",
				initialTag: "v1.0.0",
				extraTags:  []string{"foo", ""},
			},
			expectedTag: "v1.0.1",
		},
		{
			name: "test with pre-relase tag",
			setup: testRepoSetup{
				scheme:     "autotag",
				nextCommit: "#patch bump",
				initialTag: "v1.0.0",
				extraTags:  []string{"v1.0.1-pre"},
			},
			expectedTag: "v1.0.1",
		},
		{
			name: "build metadata",
			setup: testRepoSetup{
				scheme:        "autotag",
				nextCommit:    "#patch bump",
				initialTag:    "v1.0.0",
				buildMetadata: "g012345678",
			},
			expectedTag: "v1.0.1+g012345678",
		},
		{
			name: "autotag scheme, [major] bump without prefix",
			setup: testRepoSetup{
				scheme:        "autotag",
				nextCommit:    "[major] this is a big release\n\nfoo bar baz\n",
				initialTag:    "1.0.0",
				disablePrefix: true,
			},
			expectedTag: "2.0.0",
		},
		// tests for conventional commits scheme. Based on:
		// https://www.conventionalcommits.org/en/v1.0.0/#summary
		// and
		// https://www.conventionalcommits.org/en/v1.0.0/#examples
		{
			name: "conventional commits, minor bump without scope",
			setup: testRepoSetup{
				scheme:     "conventional",
				nextCommit: "feat: allow provided config object to extend other configs",
				initialTag: "v1.0.0",
			},
			expectedTag: "v1.1.0",
		},
		{
			name: "conventional commits, minor bump with scope",
			setup: testRepoSetup{
				scheme:     "conventional",
				nextCommit: "feat(lang): add polish language",
				initialTag: "v1.0.0",
			},
			expectedTag: "v1.1.0",
		},
		{
			name: "conventional commits, breaking change via ! appended to type",
			setup: testRepoSetup{
				scheme:     "conventional",
				nextCommit: "refactor!: drop support for Node 6",
				initialTag: "v1.0.0",
			},
			expectedTag: "v2.0.0",
		},
		{
			name: "conventional commits, breaking change via ! appended to type/scope",
			setup: testRepoSetup{
				scheme:     "conventional",
				nextCommit: "refactor(runtime)!: drop support for Node 6",
				initialTag: "v1.0.0",
			},
			expectedTag: "v2.0.0",
		},
		{
			name: "conventional commits, breaking change via footer",
			setup: testRepoSetup{
				scheme:     "conventional",
				nextCommit: "feat: allow provided config object to extend other configs\n\nbody before footer\n\nBREAKING CHANGE: non-backwards compatible",
				initialTag: "v1.0.0",
			},
			expectedTag: "v2.0.0",
		},
		{
			name: "conventional commits, patch/minor bump",
			setup: testRepoSetup{
				scheme:     "conventional",
				nextCommit: "fix: correct minor typos in code",
				initialTag: "v1.0.0",
			},
			expectedTag: "v1.0.1",
		},
		{
			name: "conventional commits, non-conforming fallback to patch bump",
			setup: testRepoSetup{
				scheme:     "conventional",
				nextCommit: "not a conventional commit message",
				initialTag: "v1.0.0",
			},
			expectedTag: "v1.0.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestRepo(t, tc.setup)
			defer cleanupTestRepo(t, r.repo)

			err := r.AutoTag()
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tags, err := r.repo.GetTags()
			checkFatal(t, err)
			assert.Contains(t, tags, tc.expectedTag)
		})
	}
}

func TestValidateSemVerBuildMetadata(t *testing.T) {
	tests := []struct {
		name  string
		meta  string
		valid bool
	}{
		{
			name:  "valid single-part metadata",
			meta:  "g123456",
			valid: true,
		},
		{
			name:  "valid multi-part metadata",
			meta:  "g123456.20200512",
			valid: true,
		},
		{
			name:  "invalid characters",
			meta:  "g123456,foo_bar",
			valid: false,
		},
		{
			name:  "empty string",
			meta:  "",
			valid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid := validateSemVerBuildMetadata(tc.meta)
			assert.Equal(t, tc.valid, valid)
		})
	}
}
