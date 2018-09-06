#!/bin/sh

set +x

pushd web
yarn install
yarn run build
popd
