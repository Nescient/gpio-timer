# gpio-timer
A GPIO timer for DerbyNet server (https://derbynet.org/).

This project uses golang to send timer messages to the derbynet server application.  It listens for GPIO pin changes.  Use the `build.sh` script to build the executable.  Install `docker` and use `go.sh run github.com/Nescient/gpio-timer` to run the executable.  Alternatively, install `go` to build and a new enough `glibc` for running.

## Edit your GPIO Pins
The current set up is using the Libre Computer Board ROC-RK3328-CC on pins 3, 5, 11, 13 (four lanes), and pin 8 (for the gate).  See the top of [gpio.go](gpio/gpio.go) for pin constants.  The GPIOs are set to pull up and trigger on the falling edge (which is what works with the beam sensors I tested).  Our project used Adafruit IR Break Beam sensors: https://www.adafruit.com/product/2167.

## Libre Computer Wiring
- https://hub.libre.computer/t/libre-computer-wiring-tool/40
- https://rockchip.fr/RK3328%20datasheet%20V1.2.pdf


## DerbyNet
- https://derbynet.org/builds/Installation-%20Windows.pdf
- http://drakedev.com/pinewood/
- https://derbynet.org/builds/docs/Developers-%20Timer%20Messages.pdf

<!-- RaceCoordinator, password "doyourbest" (without the quotes): this role can do anything at all. -->
<!-- RaceCrew, password "murphy" (as in Don Murphy, not Murphy’s Law): this role can do things -->
<!-- like check racers in, but not erase the database -->

## GPIO
- https://stackoverflow.com/questions/51310506/using-c-libgpiod-library-how-can-i-set-gpio-lines-to-be-outputs-and-manipulat