# gpio-timer
A GPIO timer for DebyNet server (https://derbynet.org/)

This project uses golang to send timer messages to the derbynet server application.  It listens for GPIO pin changes.  Use the `build.sh` script to build the executable.  Install `docker` and use `go.sh run github.com/Nescient/gpio-timer` to run the executable.  Alternatively, install `go` to build and a new enough `glibc` for running.

## Libre Computer Wiring
https://hub.libre.computer/t/libre-computer-wiring-tool/40
