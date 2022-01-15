#!/usr/bin/env bash

set -ex
if ! git diff --quiet; then
    echo "You have uncommitted changes, please commit them before running release."
    exit 1
fi

RE='[^0-9]*\([0-9]*\)[.]\([0-9]*\)[.]\([0-9]*\)\([0-9A-Za-z-]*\)'

step="$1"
if [ -z "$1" ]
then
  step="patch"
fi

base="$2"
if [ -z "$2" ]
then
  base=$(git tag 2>/dev/null| tail -n 1)
  if [ -z "$base" ]
  then
    base=0.0.0
  fi
fi

# shellcheck disable=SC2001
MAJOR=$(echo $base | sed -e "s#$RE#\1#")
# shellcheck disable=SC2001
MINOR=$(echo $base | sed -e "s#$RE#\2#")
# shellcheck disable=SC2001
PATCH=$(echo $base | sed -e "s#$RE#\3#")

case "$step" in
  patch)
      PATCH=$(( PATCH + 1 ))
  ;;
  minor)
      MINOR=$(( MINOR + 1 ))
      PATCH=0
  ;;
  major)
    MAJOR=$(( MAJOR + 1 ))
    MINOR=0
    PATCH=0
  ;;

esac

image_version_file="VERSION"
new_version="v$MAJOR.$MINOR.$PATCH"

echo "$new_version" > ${image_version_file}
sed -i "" "s/appVersion: .*/appVersion: $new_version/g" ./chart/Chart.yaml

git add ${image_version_file} ./chart/Chart.yaml
git commit -am "Release $new_version
[ci deploy]"
git tag "$new_version"

echo "New release $new_version ready to be pushed"
