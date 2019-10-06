const int ledPin = 9; // LED链接到9号脚

void setup()
{
    // 初始化串口通讯
    Serial.begin(9600);
    // LED作为输出
    pinMode(ledPin, OUTPUT);
}

void loop()
{
    byte brightness;

    // 检查串口数据
    if (Serial.available())
    {
        // 读取最新到达的数据:
        brightness = Serial.read();
        // 设置LED亮度
        analogWrite(ledPin, brightness);
    }
}