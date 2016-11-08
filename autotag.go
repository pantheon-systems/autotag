package autotag

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"regexp"

	"github.com/gogits/git"
	"github.com/hashicorp/go-version"
)

var (
	majorRex   = regexp.MustCompile(`(?i)\[major\]|\#major`)
	minorRex   = regexp.MustCompile(`(?i)\[minor\]|\#minor`)
	patchRex   = regexp.MustCompile(`(?i)\[patch\]|\#patch`)
	versionRex = regexp.MustCompile(`^v([\d]+\.?.*)`)
)

// GitRepo represents a repository we want to run actions against
type GitRepo struct {
	repo *git.Repository

	currentVersion *version.Version
	currentTag     *git.Commit
	newVersion     *version.Version
	branch         string
	branchID       string // commit id of the branch latest commit (where we will apply the tag)
}

// NewRepo is a constructor for a repo object, parsing the tags that exist
func NewRepo(repoPath, branch string) (*GitRepo, error) {
	if branch == "" {
		return nil, fmt.Errorf("must specify a branch")
	}

	log.Println("Opening repo at ", repoPath+"/.git")
	repo, err := git.OpenRepository(repoPath + "/.git")
	if err != nil {
		return nil, err
	}

	r := &GitRepo{
		repo:   repo,
		branch: branch,
	}

	err = r.parseTags()
	if err != nil {
		return nil, err
	}

	if err := r.calcVersion(); err != nil {
		return nil, err
	}

	return r, nil
}

// Temp shim that uses git to parse the tags on the repo.
// this is because the library at the moment does not parse packed refs
// TODO: move to https://github.com/src-d/go-git as the pure-go library. It supports everything we need
func (r *GitRepo) getTags() (map[string]string, error) {
	tags := make(map[string]string)
	var outb, errb bytes.Buffer

	gitbin, err := exec.LookPath("git")
	if err != nil {
		return tags, fmt.Errorf("git executable not found: %s", err)
	}

	p := r.repo.Path
	if strings.Contains(p, "/.git") {
		p, err = filepath.Abs(p + "/../")
		if err != nil {
			return tags, err
		}
	}

	cmd := exec.Command(gitbin, "show-ref", "--tags")
	cmd.Dir = p

	cmd.Stderr = &errb
	cmd.Stdout = &outb
	err = cmd.Run()
	if err != nil {
		return tags, fmt.Errorf("failed listing tags '%s': %s", errb.String(), err)
	}

	scanner := bufio.NewScanner(&outb)
	for scanner.Scan() {
		t := strings.Split(scanner.Text(), " refs/tags/")
		tags[t[1]] = t[0]
	}

	return tags, nil
}

// Parse tags on repo, sort them, and store the most recent revision in the repo object
func (r *GitRepo) parseTags() error {
	log.Println("Parsing repository tags")

	versions := make(map[*version.Version]*git.Commit)

	tags, err := r.getTags()
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %s", err.Error())
	}

	for tag, commit := range tags {
		v, err := maybeVersionFromTag(tag)
		if err != nil {
			log.Println("skipping non version tag: ", tag)
			continue
		}

		if v == nil {
			log.Println("skipping non version tag: ", tag)
			continue
		}

		c, err := r.repo.GetCommit(commit)
		if err != nil {
			return fmt.Errorf("error reading commit '%s':  %s", commit, err)
		}
		versions[v] = c
	}

	keys := make([]*version.Version, 0, len(versions))
	for key := range versions {
		keys = append(keys, key)
	}
	sort.Sort(sort.Reverse(version.Collection(keys)))

	// set the current versions
	if len(keys) >= 1 {
		v := keys[0]
		r.currentVersion = v
		r.currentTag = versions[v]

		//		log.Printf("Current latest version is %s at obj: %s id: %s", r.currentVersion, r.currentTag.Object, r.currentTag.Id)
		return nil
	}

	return fmt.Errorf("no version tags found")

}

func maybeVersionFromTag(tag string) (*version.Version, error) {
	if tag == "" {
		return nil, fmt.Errorf("empty tag not supported")
	}

	ver, vErr := parseVersion(tag)
	if vErr != nil {
		return nil, fmt.Errorf("couldn't parse version %s: %s", tag, vErr)
	}
	return ver, nil
}

// parseVersion returns a version object from a parsed string. This normalizes semver strings, and adds the ability to parse strings with 'v' leader. so that `v1.0.1`->     `1.0.1`  which we need for berkshelf to work
func parseVersion(v string) (*version.Version, error) {
	if versionRex.MatchString(v) {
		m := versionRex.FindStringSubmatch(v)
		if len(m) >= 2 {
			v = m[1]
		}
	}

	nVersion, err := version.NewVersion(v)
	if err != nil && nVersion != nil && len(nVersion.Segments()) >= 1 {
		return nVersion, err
	}
	return nVersion, nil
}

// LatestVersion Reports the Lattest version of the given repo
// TODO:(jnelson) this could be more intelligent, looking for a nil new and reporitng the latest version found if we refactor autobump at some point Mon Sep 14 13:05:49 2015
func (r *GitRepo) LatestVersion() string {
	return fmt.Sprintf("%s", r.newVersion)
}

func (r *GitRepo) retrieveBranchInfo() error {
	id, err := r.repo.GetCommitIdOfBranch(r.branch)
	if err != nil {
		return fmt.Errorf("error getting head commit: %s ", err.Error())
	}

	r.branchID = id
	return nil
}

// calcVersion looks over commits since the last tag, and will apply the version bump needed. It will patch if no other instruction is found
// it populates the repo.newVersion with the new calculated version
func (r *GitRepo) calcVersion() error {
	r.newVersion = r.currentVersion
	if err := r.retrieveBranchInfo(); err != nil {
		return err
	}

	l, err := r.repo.CommitsBefore(r.branchID)
	if err != nil {
		log.Printf("Error loading history for tag '%s': %s ", r.currentVersion, err.Error())
	}
	log.Printf("Checking commits from %s to %s ", r.branchID, r.currentTag.Id)

	// Sort the commits oldest to newest. Then process each commit for bumper commands.
	start := false
	for e := l.Back(); e != nil; e = e.Prev() {
		commit := e.Value.(*git.Commit)
		if commit == nil {
			return fmt.Errorf("commit pointed to nil object. This should not happen.\n")
		}

		// we scan from the first commit till the tagCommit.
		if commit.Id == r.currentTag.Id {
			start = true
			continue
		}

		// unless we have found the commit theres no need to process
		if !start {
			continue
		}

		v, nerr := r.parseCommit(commit)
		if nerr != nil {
			log.Fatal(err)
		}

		if v != nil {
			r.newVersion = v
		}

	}

	// if there is no movement on the version from commits, bump patch
	if r.newVersion == r.currentVersion {
		if r.newVersion, err = patchBumper.bump(r.currentVersion); err != nil {
			return err
		}
	}
	return nil
}

// AutoTag applies the new version tag thats calculated
func (r *GitRepo) AutoTag() error {
	return r.tagNewVersion()
}

func (r *GitRepo) tagNewVersion() error {
	// TODO:(jnelson) These should be configurable? Mon Sep 14 12:02:52 2015
	tagName := fmt.Sprintf("v%s", r.newVersion.String())

	log.Println("Writing Tag", tagName)
	err := r.repo.CreateTag(tagName, r.branchID)
	if err != nil {
		return fmt.Errorf("error creating tag: %s", err.Error())
	}
	return nil
}

// parseLog looks at HEAD commit see if we want to increment major/minor/patch
func (r *GitRepo) parseCommit(commit *git.Commit) (*version.Version, error) {
	var b bumper
	msg := commit.Message()
	log.Printf("Parsing %s: %s\n", commit.Id, msg)

	if majorRex.MatchString(msg) {
		log.Println("major bump")
		b = majorBumper
	}

	if minorRex.MatchString(msg) {
		log.Println("minor bump")
		b = minorBumper
	}

	if patchRex.MatchString(msg) {
		log.Println("patch bump")
		b = patchBumper
	}

	if b != nil {
		return b.bump(r.currentVersion)
	}

	return nil, nil
}

// MajorBump will bump the version one major rev 1.0.0 -> 2.0.0
func (r *GitRepo) MajorBump() (*version.Version, error) {
	return majorBumper.bump(r.currentVersion)
}

// MinorBump will bump the version one minor rev 1.1.0 -> 1.2.0
func (r *GitRepo) MinorBump() (*version.Version, error) {
	return minorBumper.bump(r.currentVersion)
}

// PatchBump will bump the version one patch rev 1.1.1 -> 1.1.2
func (r *GitRepo) PatchBump() (*version.Version, error) {
	return patchBumper.bump(r.currentVersion)
}
