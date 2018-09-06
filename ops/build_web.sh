#!/bin/sh

set +x

pushd web
npm install
npm run build
popd
