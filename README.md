# gpio-timer
A GPIO timer for DerbyNet server (https://derbynet.org/).

This project uses golang to send timer messages to the derbynet server application.  It listens for GPIO pin changes.  Use the `build.sh` script to build the executable.  Install `docker` and use `go.sh run github.com/Nescient/gpio-timer` to run the executable.  Alternatively, install `go` to build and a new enough `glibc` for running.

## Libre Computer Wiring
https://hub.libre.computer/t/libre-computer-wiring-tool/40


## DerbyNet
- https://derbynet.org/builds/Installation-%20Windows.pdf
- http://drakedev.com/pinewood/
- https://derbynet.org/builds/docs/Developers-%20Timer%20Messages.pdf

<!-- RaceCoordinator, password "doyourbest" (without the quotes): this role can do anything at all. -->
<!-- RaceCrew, password "murphy" (as in Don Murphy, not Murphy’s Law): this role can do things -->
<!-- like check racers in, but not erase the database -->

## GPIO
- https://stackoverflow.com/questions/51310506/using-c-libgpiod-library-how-can-i-set-gpio-lines-to-be-outputs-and-manipulat