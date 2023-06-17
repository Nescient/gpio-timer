# gpio-timer
A GPIO timer for DerbyNet server (https://derbynet.org/).

- ![#f03c15](https://placehold.co/15x15/f03c15/f03c15.png) `WORK IN PROGRESS, NO RELEASE`
```diff
- DO NOT USE THIS PROJECT YET.  I PLAN TO HAVE IT DONE BY JULY 1, 2023. PLEASE CHECK BACK THEN.
```
- ![#f03c15](https://placehold.co/15x15/f03c15/f03c15.png) `WORK IN PROGRESS, NO RELEASE`

This project uses golang to send timer messages to the derbynet server application.  It listens for GPIO pin changes.  Use the `build.sh` script to build the executable.  Install `docker` and use `go.sh run github.com/Nescient/gpio-timer` to run the executable.  Alternatively, install `go` to build and a new enough `glibc` for running.

## Libre Computer Wiring
https://hub.libre.computer/t/libre-computer-wiring-tool/40


## DerbyNet
- https://derbynet.org/builds/Installation-%20Windows.pdf
- http://drakedev.com/pinewood/
- https://derbynet.org/builds/docs/Developers-%20Timer%20Messages.pdf

## GPIO
- https://stackoverflow.com/questions/51310506/using-c-libgpiod-library-how-can-i-set-gpio-lines-to-be-outputs-and-manipulat