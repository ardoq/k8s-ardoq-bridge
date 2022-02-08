#!/usr/bin/env bash

#set -ex
if ! git diff --quiet; then
    echo "You have uncommitted changes, please commit them before running release."
    exit 1
fi

HELM_DIR=helm/chart
GH_OWNER=ardoq
HELM_REP=k8s-ardoq-bridge

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

sed -i "" "s/version: .*/version: $new_version/g" $HELM_DIR/Chart.yaml
sed -i "" "s/appVersion: .*/appVersion: $latest_tagged_version/g" $HELM_DIR/Chart.yaml

git add CHART_VERSION $HELM_DIR
git commit -m "Upgraded helm chart version to $new_version appVersion to $latest_tagged_version"


echo "Upgraded helm chart version to $new_version appVersion to $latest_tagged_version ready to be pushed"

echo "Releasing Helm chart"

function setup_chart_releaser() {
    curl -OL https://github.com/helm/chart-releaser/releases/download/v1.3.0/chart-releaser_1.3.0_linux_amd64.tar.gz
    echo "baed2315a9bb799efb71d512c5198a2a3b8dcd139d7f22f878777cffcd649a37  chart-releaser_1.3.0_linux_amd64.tar.gz" | sha256sum -c -
    tar xzvf chart-releaser_1.3.0_linux_amd64.tar.gz cr
}
function cleanup() {
    rm cr
}

echo "Linting"
helm lint $HELM_DIR

echo "setting up chart releaser"
setup_chart_releaser

echo "package helm chart"
./cr package $HELM_DIR -p helm

echo "Uplocad helm chart"
./cr upload -o $GH_OWNER -r $HELM_REP --skip-existing -p helm
#git add helm/k8s-ardoq-bridge-*
#git commit -m '[automated commit] uploaded archived helm chart'

echo "Index Helm chart"
./cr index -o $GH_OWNER -r $HELM_REP -c https://raw.githubusercontent.com/$GH_OWNER/$HELM_REP/main/ -i helm/index.yaml -p helm --push
#git add helm/index.yaml
#git commit -m '[automated commit] uploaded index file'

echo "Push all staged changes"
#git push

echo "cleanup"
cleanup