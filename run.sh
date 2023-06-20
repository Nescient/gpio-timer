#!/bin/sh

./go.sh generate
# ./go.sh run github.com/Nescient/gpio-timer

docker run --net=host --rm -v "$PWD":/usr/src/myapp \
   -w /usr/src/myapp -e HOME=/usr/src/myapp --device=/dev/gpiochip1 \
   --device=/dev/gpiochip2 golang:1.20 go run github.com/Nescient/gpio-timer

