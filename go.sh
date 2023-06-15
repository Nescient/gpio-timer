#!/bin/sh

docker run --user $(id -u):$(id -g) --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.20 go $@
