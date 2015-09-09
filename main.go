package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/jessevdk/go-flags"
	"github.com/libgit2/git2go"
)

type Options struct {
	Null bool `short:"n" description:"Enable dry run mode"`
	Auto bool `short:"n" description:"Enable dry run mode"`
}

var Opts Options

func init() {
	_, err := flags.Parse(&Opts)
	if err != nil {
		os.Exit(1)
	}
}

type GitRepo struct {
	*git.Repository
}

func NewRepo() *GitRepo {
	repoPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	return &GitRepo{repo}
}

func init() {
	_, err := flags.Parse(&Opts)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	r := NewRepo()
	r.parseTags()
	os.Exit(1)
}

type Versions []*version.Version

func NewSortedVersionCollection(vs Versions) Versions {
	sort.Sort(version.Collection(vs))
	return vs
}

func (r *GitRepo) parseTags() {
	iter, err := r.NewReferenceIterator()
	if err != nil {
		log.Fatal(err)
	}

	versions := make(Versions, 10)
	for ref, err := iter.Next(); err == nil; {
		if v, err := maybeVersionFromTag(ref); err != nil {
			versions = append(versions, v)
		}
	}

	versions = NewSortedVersionCollection(versions)
	latest := versions[len(versions)-1]
	fmt.Println("Latest version is", latest)

	s := latest.Segments()
	nextMajor := s[0] + 1
	nextMinor := s[1] + 1
	nextBump := s[2] + 1

	if Opts.Null == true {
		fmt.Printf("Next major: %d.%d.%d\n", nextMajor, s[1], s[2])
		fmt.Printf("Next minor: %d.%d.%d\n", s[0], nextMinor, s[2])
		fmt.Printf("Next patch: %d.%d.%d\n", s[0], s[1], nextBump)
		os.Exit(0)
	}

}

func maybeVersionFromTag(ref *git.Reference) (*version.Version, error) {
	if ref.IsTag() {
		name := ref.Shorthand()
		//			fmt.Println("Found Tag: ", name)
		ver, vErr := version.NewVersion(name)
		if vErr != nil {
			return ver, fmt.Errorf("couldn't parse version %s: %s", name, vErr)
		}
		return ver, nil
	}
	return nil, nil
}

/*
func (c *Candidate) parseMessage() {
	   logs = ShellUtils.sh "git log --abbrev-commit --format=oneline #{last_tag}.."
	   guess = if logs =~ /\[major\]|\#major/i
	             :major
	           elsif logs =~ /\[minor\]|\#minor/i
	             :minor
	           elsif logs =~ /\[prerelease\s?(#{Prerelease::TYPE_FORMAT})?\]|\#prerelease\-?(#{Prerelease::TYPE_FORMAT})?/
	             prerelease_type = $1 || $2
	             :prerelease
	           elsif logs =~ /\[patch\]|\#patch/i
	             :patch
	           else
	             options[:default] or :build
	           end
	   bump!(guess, prerelease_type: prerelease_type)
}

*/
