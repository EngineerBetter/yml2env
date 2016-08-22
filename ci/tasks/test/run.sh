#!/bin/bash
set -xe

export PATH=$PATH:$PWD
export GOPATH=$PWD/gopath

cd ${GOPATH}/src/github.com/EngineerBetter/yml2env
ginkgo -r -v