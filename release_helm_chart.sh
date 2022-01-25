#!/usr/bin/env bash

#set -ex
if ! git diff --quiet; then
    echo "You have uncommitted changes, please commit them before running release."
    exit 1
fi
image_version_file="VERSION"
latest_tagged_version="$1"
if [ -z "$1" ]
then
  latest_tagged_version=$(git tag --sort=v:refname 2>/dev/null| tail -n 1)
  if [ -z "$latest_tagged_version" ]
  then
    latest_tagged_version=$(cat ${image_version_file})
  fi
fi

sed -i "" "s/appVersion: .*/appVersion: $latest_tagged_version/g" ./chart/Chart.yaml

git add ./chart/
git commit -am "Upgraded helm chart appVersion to $latest_tagged_version"


echo "Upgraded helm chart appVersion to $latest_tagged_version ready to be pushed"
