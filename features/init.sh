#!/bin/bash

if [ ! -z "${TRAVIS+x}" ];
then
    git config --global user.name "whatever";
    git config --global user.email "whatever@example.com";
fi

# Configure name

# Init
rm -rf test > /dev/null;
git init test;

cd test || exit 1;

git config --local user.name "whatever";
git config --local user.email "whatever@example.com";
