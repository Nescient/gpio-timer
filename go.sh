#!/bin/sh

docker run --net=host --user $(id -u):$(id -g) --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e HOME=/usr/src/myapp golang:1.24 go $@
