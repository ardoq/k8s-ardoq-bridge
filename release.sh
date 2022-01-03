#!/usr/bin/env bash

set -ex
if ! git diff --quiet; then
    echo "You have uncommitted changes, please commit them before running release."
    exit 1
fi
image_version_file="image_version.txt"
new_version=$(( $(cat ${image_version_file}) + 1 ))

echo "$new_version" > ${image_version_file}

git add ${image_version_file}
git commit -am "Release v$new_version
[ci deploy]"
git tag "v$new_version"

echo "New release $new_version ready to be pushed"
