package autotag

import (
	"fmt"
	"log"
	"sort"
	"time"

	"regexp"

	"github.com/hashicorp/go-version"
	"github.com/libgit2/git2go"
)

var (
	majorRex   = regexp.MustCompile(`(?i)\[major\]|\#major`)
	minorRex   = regexp.MustCompile(`(?i)\[minor\]|\#minor`)
	patchRex   = regexp.MustCompile(`(?i)\[patch\]|\#patch`)
	versionRex = regexp.MustCompile(`^v([\d]+\.?.*)`)
)

type GitRepo struct {
	Repo           *git.Repository
	walker         *git.RevWalk
	currentVersion *version.Version
	newVersion     *version.Version
}

// NewRepo is a constructor for a repo object, parsing the tags that exist
// The caller is responsible for calling Free() on the embedded Repo when they are done with it
func NewRepo(repoPath string) (*GitRepo, error) {

	log.Println("Opening repo at ", repoPath)
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return nil, err
	}

	w, err := repo.Walk()
	if err != nil {
		return nil, err
	}

	r := &GitRepo{
		Repo:   repo,
		walker: w,
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

func (r *GitRepo) parseTags() error {
	log.Print("Parsing repo tags")
	iter, err := r.Repo.NewReferenceIterator()
	if err != nil {
		return fmt.Errorf("ouch, couldn't get a ref on repo this shouldn't happen ever: %s", err)
	}

	var versions []*version.Version
	var lastParsedTagRef string

	// step through the refrences
	for ref, err := iter.Next(); err == nil; ref, err = iter.Next() {
		log.Println("checking ref: ", ref.Name())
		if v, err := maybeVersionFromTag(ref); err == nil {
			versions = append(versions, v)
			lastParsedTagRef = ref.Target().String()
		}
	}

	cRange := fmt.Sprintf("%s..HEAD", lastParsedTagRef)
	//	log.Println("Check range:", cRange)
	if wErr := r.walker.PushRange(cRange); wErr != nil {
		log.Fatal(wErr)
	}

	sort.Sort(version.Collection(versions))

	// if more than one version, take the last one
	if itemLen := len(versions); itemLen >= 1 {
		r.currentVersion = versions[itemLen-1]
	}

	return nil
}

func maybeVersionFromTag(ref *git.Reference) (*version.Version, error) {
	if ref == nil {
		return nil, fmt.Errorf("nil refrence")
	}

	if ref.IsTag() {
		name := ref.Shorthand()
		log.Println("Found Tag: ", name)
		ver, vErr := parseVersion(name)
		if vErr != nil {
			return nil, fmt.Errorf("couldn't parse version %s: %s", name, vErr)
		}
		return ver, nil
	}
	return nil, fmt.Errorf("not a tag")
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
	if err != nil && len(nVersion.Segments()) >= 1 {
		return nVersion, err
	}
	return nVersion, nil
}

// Report the Lattest version
// TODO:(jnelson) this could be more intelligent, looking for a nil new and reporitng the latest version found if we refactor autobump at some point Mon Sep 14 13:05:49 2015
func (r *GitRepo) LatestVersion() string {
	return fmt.Sprintf("v%s", r.newVersion)
}

// calcVersion looks over commits since the last tag, and will apply the version bump needed. It will patch if no other instruction is found
// it populates the repo.newVersion with the new calculated version
func (r *GitRepo) calcVersion() error {
	w := r.walker
	defer w.Free()

	r.newVersion = r.currentVersion

	err := w.Iterate(func(c *git.Commit) bool {
		v, err := r.parseCommit(c)
		if err != nil {
			log.Fatal(err)
		}
		if v != nil {
			r.newVersion = v
		}

		return true
	})

	if err != nil {
		log.Fatal(err)
	}

	// if there is no movement on the version from commits, bump patch
	if r.newVersion == r.currentVersion {
		if r.newVersion, err = patchBumper.bump(r.currentVersion); err != nil {
			return err
		}
	}
	return nil
}

// AutoBump applies the new version tag thats calculated
func (r *GitRepo) AutoTag() error {
	if err := r.tagNewVersion(); err != nil {
		return err
	}

	return nil
}

func (r *GitRepo) tagNewVersion() error {
	// TODO:(jnelson) These should be configurable? Mon Sep 14 12:02:52 2015
	sig := &git.Signature{
		Name:  "AutoTag",
		Email: "noreply",
		When:  time.Now(),
	}

	currentBranch, err := r.Repo.Head()
	if err != nil {
		return err
	}

	tip, err := r.Repo.LookupCommit(currentBranch.Target())
	if err != nil {
		return err
	}

	tagName := fmt.Sprintf("v%s", r.newVersion.String())
	if _, err := r.Repo.Tags.Create(tagName, tip, sig, fmt.Sprintf("AutoTag to %s", tagName)); err != nil {
		return err
	}
	return nil
}

// parseLog looks at HEAD commit see if we want to increment major/minor/patch
func (r *GitRepo) parseCommit(commit *git.Commit) (*version.Version, error) {
	var b bumper
	msg := commit.Message()
	log.Println("Parsing ", msg, ":")

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
