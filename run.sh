#!/bin/sh

./go.sh generate
# ./go.sh run github.com/Nescient/gpio-timer

docker stop derby-timer
docker rm derby-timer

docker run -itd --name derby-timer --restart unless-stopped --net=host \
   -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e HOME=/usr/src/myapp \
   --device=/dev/gpiochip1 --device=/dev/gpiochip2 --device=/dev/gpiochip3 --device=/dev/gpiochip4 \
   golang:1.24 go run github.com/Nescient/gpio-timer

docker logs --follow derby-timer
