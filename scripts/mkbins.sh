#!/bin/bash

set -xe

if [ ! -d ~/bin ] ; then
  mkdir -p ~/bin
fi

if [ ! -f ~/bin/gh-release ] ; then
  curl -L https://github.com/progrium/gh-release/releases/download/v2.2.0/gh-release_2.2.0_linux_x86_64.tgz  | tar -xzv
  mv gh-release ~/bin/gh-release
  chmod 755 ~/bin/gh-release
fi
