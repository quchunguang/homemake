 // Dimmer - 用串口发送数据
 
 import processing.serial.*;
 Serial port;
 
 void setup() {
 size(256, 150);
 
 print("Available serial ports:");
 println(Serial.list());
 
 // 使用列表里的第一个可用串口. 选择和arduino对应的串口和通信速率
 port = new Serial(this, Serial.list()[0], 9600);  
 
 // 如果你知道arduino的串口，就直接这样写
 //port = new Serial(this, "COM1", 9600);
 }
 
 void draw() {
 // 画一个由黑到白的渐变图
 for (int i = 0; i < 256; i++) {
 stroke(i);
 line(i, 0, i, 150);
 }
 
 // 以单字节形式把鼠标x坐标信息发送到串口
 
 port.write(mouseX);
 }
