#!/bin/bash

set -xe

if [  -d "$HOME/go1.5" ] ; then
  echo "go 1.5 installed skipping"
  exit 0
fi


curl -o "$HOME/go1.5.tar.gz" https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz
tar -C "$HOME/" -xzvf "$HOME/go1.5.tar.gz"
mv "$HOME/go" "$HOME/go1.5"
