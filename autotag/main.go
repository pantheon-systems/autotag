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
	JustVersion bool   `short:"n" description:"Just output the next version, don't autotag"`
	Verbose     bool   `short:"v" description:"Enable verbose logging"`
	Branch      string `short:"b" long:"branch" description:"Git branch to scan" default:"master" `
	RepoPath    string `short:"r" long:"repo" description:"Path to the repo" default:"./" `
}

var opts Options

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	log.SetOutput(ioutil.Discard)
	if opts.Verbose {
		log.SetOutput(os.Stderr)
	}

	r, err := autotag.NewRepo(opts.RepoPath, opts.Branch)
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
