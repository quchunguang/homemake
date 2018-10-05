MyFirstSketch
=============

Demo of how to work with arduino-cli. [arduino-cli](https://github.com/arduino/arduino-cli).
No need to install the offical Arduino IDE!

```bash
# or download latest release
go get -u github.com/arduino/arduino-cli
# root for everything
mkdir $HOME/Arduino

# update averiable cores
arduino-cli core update-index
# without the core installed, shows the id and port only
arduino-cli board list
# search for the core for the board
arduino-cli core search arduino
arduino-cli core search uno
# install core
arduino-cli core install arduino:avr
# ensure the core has installed successfully 
arduino-cli core list
# BUG: should shows the detail of board connected
arduino-cli board list
# shows the USB device belongs to dialout group
ls /dev/USB0
# make USB port accessable by $USER
sudo usermod -a -G dialout $USER
# MUST reboot to enable the change
sudo reboot

# look at the name of a library
arduino-cli lib search wifi101
# install the library by name
arduino-cli lib install "WiFi101"

# create the project folder, and then edit the source in MyFirstSketch/MyFirstSketch.ino
arduino-cli sketch new MyFirstSketch
cd $HOME/Arduino/MyFirstSketch
# compile
arduino-cli compile --fqbn arduino:avr:uno .
# upload
arduino-cli upload -p /dev/ttyUSB0 --fqbn arduino:avr:uno .

# altnatively, `make`
make
```