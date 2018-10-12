#!/usr/bin/env python3
# tested on rpi3 model B
import RPi.GPIO as GPIO
import time
# print(GPIO.RPI_INFO)

pin_ground = 6
pin_led = 8

GPIO.setmode(GPIO.BOARD) # pin number
GPIO.setup(pin_led, GPIO.OUT, initial=0)
try:
    while True:
        GPIO.output(pin_led, GPIO.HIGH)
        print('1', end='', flush=True)
        time.sleep(.5)

        GPIO.output(pin_led, GPIO.LOW)
        print('0', end='', flush=True)
        time.sleep(.5)
except KeyboardInterrupt:
    GPIO.cleanup()
    print('Exiting...')
