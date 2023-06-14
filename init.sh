#!/bin/sh

conan install . -s build_type=Debug -b missing || exit $?
cmake --preset=dev || exit $?
cmake --build --preset=dev || exit $?
ctest --preset=dev