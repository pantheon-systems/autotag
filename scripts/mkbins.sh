#!/bin/bash

set -xe

if [ ! -d ~/bin ] ; then
  mkdir -p ~/bin
fi

if [ ! -f ~/bin/gh-release ] ; then
  # when https://github.com/progrium/gh-release/pull/13 is merged we can go back
  # to upstream on this for now I am using my build from that pr -jesse
  # curl -L https://github.com/progrium/gh-release/releases/download/v2.2.0/gh-release_2.2.0_linux_x86_64.tgz  | tar -xzv
  # mv gh-release ~/bin/gh-release
  curl -L https://www.dropbox.com/s/4k3eq7xpehwwqr5/gh-release?dl=0 -o ~/bin/gh-release
  chmod 755 ~/bin/gh-release
fi
