#!/bin/bash
# README.md for detail
arduino-cli compile --fqbn arduino:avr:uno .
arduino-cli upload -p /dev/ttyUSB0 --fqbn arduino:avr:uno .
