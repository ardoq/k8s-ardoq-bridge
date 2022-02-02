#!/usr/bin/env bash

#set -ex
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

chart_version_file="CHART_VERSION"

base="$2"
if [ -z "$2" ]
then
  base=$(cat ${chart_version_file})
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

new_version="$MAJOR.$MINOR.$PATCH"

echo "$new_version" > ${chart_version_file}

image_version_file="VERSION"
latest_tagged_version="$3"
if [ -z "$3" ]
then
  latest_tagged_version=$(git tag --sort=v:refname 2>/dev/null| tail -n 1)
  if [ -z "$latest_tagged_version" ]
  then
    latest_tagged_version=$(cat ${image_version_file})
  fi
fi

sed -i "" "s/version: .*/version: $new_version/g" ./helm/chart/Chart.yaml
sed -i "" "s/appVersion: .*/appVersion: $latest_tagged_version/g" ./helm/chart/Chart.yaml

git add CHART_VERSION ./helm/chart/
git commit -am "Upgraded helm chart version to $new_version appVersion to $latest_tagged_version"


echo "Upgraded helm chart version to $new_version appVersion to $latest_tagged_version ready to be pushed"
