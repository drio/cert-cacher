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

gh release create v$new_version --notes "v$new_version" *.linux.amd64 *.darwin.arm64 *.windows.amd64

notify "New release done 🚀"
