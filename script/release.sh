#!/bin/sh

set -e
latest_tag=$(git describe --abbrev=0 --tags)
goxc
ghr -u gainings -r gainings $latest_tag dist/snapshot/
