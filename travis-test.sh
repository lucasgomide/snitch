#!/usr/bin/env bash
set -e

for d in $(go list ./... | grep -v vendor); do
  go test -coverprofile=profile.out $d
  if [ -f profile.out ]; then
    cat profile.out > `echo $d | awk -F "/" '{print $NF}'`.coverprofile
    rm profile.out
  fi
done
gover && goveralls -coverprofile=gover.coverprofile -service travis-ci
