#!/bin/bash

set -eu

current=$(cat version.go | tail -1 | head -1 | awk -F= '{print $2}')
echo "Current version is:"
echo $current
echo -ne "Enter new version: "
read new_version

(
	cat <<EOF
package main

const version = "$new_version"
EOF
) >version.go

git add version.go
git commit -m 'bump version'
git tag "v$new_version"
git push origin
git push origin v$new_version

make clean build

gh release create v$new_version --notes "v$new_version" ./cert-cacher.arm64.osx ./cert-cacher.amd64.linux ./cert-cacher.amd64.windows

notify "New release done 🚀"
