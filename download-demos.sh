#!/bin/bash

set -e

mkdir -p cs-demos

curl -s -L -o /tmp/demos.txt https://gitlab.com/akiver/cs-demos/-/raw/main/demos.txt
files=$(cat /tmp/demos.txt)

for file in $files; do
  demoPath="cs-demos/$file"
  if [ -f $demoPath ]; then
    continue
  fi

  echo "Downloading demo $demoPath"
  curl -s -L -o $demoPath https://gitlab.com/akiver/cs-demos/-/raw/main/$file --create-dirs
done

echo 'Done'
