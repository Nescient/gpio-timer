#!/bin/sh

./go.sh generate
# ./go.sh run github.com/Nescient/gpio-timer

docker run -itd --name derby-timer --restart unless-stopped --net=host \
   -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e HOME=/usr/src/myapp \
   --device=/dev/gpiochip1 --device=/dev/gpiochip2 --device=/dev/gpiochip3 --device=/dev/gpiochip4 \
   golang:1.20 go run github.com/Nescient/gpio-timer

echo "To see output in the terminal, run"
echo "docker logs derby-timer"
