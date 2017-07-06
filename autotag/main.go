package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/pantheon-systems/autotag"
)

// Options holds the CLI args
type Options struct {
	JustVersion         bool   `short:"n" description:"Just output the next version, don't autotag"`
	Verbose             bool   `short:"v" description:"Enable verbose logging"`
	Branch              string `short:"b" long:"branch" description:"Git branch to scan" default:"master" `
	RepoPath            string `short:"r" long:"repo" description:"Path to the repo" default:"./" `
	PreReleaseName      string `short:"p" long:"pre-release-name" description:"create a pre-release tag with this name (can be: alpha|beta|pre|rc)"`
	PreReleaseTimestamp string `short:"T" long:"pre-release-timestamp" description:"create a pre-release tag and append a timestamp (can be: datetime|epoch)"`
}

var opts Options

const (
	// epochTsLayout is the UNIX epoch time format
	epochTsLayout = "epoch"

	// datetimeTsLayout is the YYYYMMDDHHMMSS time format
	datetimeTsLayout = "20060102150405"
)

func timestampLayoutFromOpts() string {
	switch opts.PreReleaseTimestamp {
	case "epoch":
		return epochTsLayout
	case "datetime":
		return datetimeTsLayout
	default:
		return ""
	}
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if err := validateOpts(); err != nil {
		log.SetOutput(os.Stderr)
		log.Fatalf("error validating flags: %s\n", err.Error())
	}
}

func validateOpts() error {
	switch opts.PreReleaseName {
	case "", "alpha", "beta", "pre", "rc":
		// nothing -- valid values
	default:
		return fmt.Errorf("-p/--pre-release-name was %q; want (alpha|beta|pre|rc)", opts.PreReleaseName)
	}

	switch opts.PreReleaseTimestamp {
	case "", "datetime", "epoch":
		// nothing -- valid values
	default:
		return fmt.Errorf("-T/--pre-release-timestamp was %q; want (datetime|epoch)", opts.PreReleaseTimestamp)
	}

	return nil
}

func main() {
	log.SetOutput(ioutil.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}

	r, err := autotag.NewRepo(autotag.GitRepoConfig{
		RepoPath:                  opts.RepoPath,
		Branch:                    opts.Branch,
		PreReleaseName:            opts.PreReleaseName,
		PreReleaseTimestampLayout: timestampLayoutFromOpts(),
	})

	if err != nil {
		fmt.Println("Error initializing: ", err)
		os.Exit(1)
	}

	// Tag unless asked otherwise
	if !opts.JustVersion {
		err = r.AutoTag()
		if err != nil {
			fmt.Println("Error auto updating version: ", err.Error())
			os.Exit(1)
		}
	}

	fmt.Println(r.LatestVersion())

	// TODO:(jnelson) Add -major -minor -patch flags for force bumps Fri Sep 11 10:04:20 2015
	os.Exit(0)
}
