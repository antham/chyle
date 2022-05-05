#!/usr/bin/env bash

rootPath=$(pwd)

cd "$rootPath/cmd" && sh "$rootPath/features/init.sh"
cd "$rootPath/cmd" && sh "$rootPath/features/merge-commits.sh"
cd "$rootPath/chyle" && sh "$rootPath/features/init.sh"
cd "$rootPath/chyle" && sh "$rootPath/features/merge-commits.sh"
cd "$rootPath/chyle/git" && sh "$rootPath/features/init.sh"
cd "$rootPath/chyle/git" && sh "$rootPath/features/merge-commits.sh"
