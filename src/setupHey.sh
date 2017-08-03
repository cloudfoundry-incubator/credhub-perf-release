#!/bin/bash

cd .. # takes us to credhub-performance/
git submodule init
git submodule update

go build -o src/credhub_cannon/hey github.com/rakyll/hey 
