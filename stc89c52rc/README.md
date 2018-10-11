stc89c52rc
==========

It demos the compile and upload serial small chips under Ubuntu.
* [入门】Linux上搭建51单片机开发环境（一） —— 环境搭建](https://blog.csdn.net/qq_21460229/article/details/73457996)
* [stcgal](https://github.com/grigorig/stcgal)
* [stc-header](https://github.com/znhocn/stc-header)


```bash
sudo apt-get install sdcc
#sudo apt install python3-pip
pip3 install stcgal # add $HOME/.local/bin to $PATH at $HOME/.profile
cd $HERE
git clone https://github.com/znhocn/stc-header.git
code main.c # using "stc-header/STC89xx.h" instead if standard c51 header file

sdcc main.c
packihx main.ihx >main.hex
stcgal -P stc89 main.hex

# If using `.bin` instead of `.hex` for upload, following two methods are equal
# objcopy -Iihex -Obinary main.hex main.bin
# makebin -p < main.hex > main.bin
```