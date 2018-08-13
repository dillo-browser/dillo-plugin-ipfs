#!/bin/sh
set -e
git update-server-info
hash="$(ipfs add -prwQH $(git ls-files) .git)"
ipfs name publish --key=dillo-ipfs /ipfs/"$hash"
