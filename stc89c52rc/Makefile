all:
	sdcc main.c
	packihx main.ihx >main.hex
	stcgal -P stc89 main.hex

clean:
	rm -f *.bin *.hex *.lk *.map *.rel *.sym *.asm *.ihx *.lst *.mem *.rst