#!/bin/bash

set -ex

g2g_path="$GOPATH/src/github.com/libgit2/git2go"
if [ -n "$CIRCLECI" ] ; then
  g2g_path="/home/ubuntu/.go_workspace/src/github.com/libgit2/git2go"
fi

cd $g2g_path
git checkout next
git submodule update --init
make install
