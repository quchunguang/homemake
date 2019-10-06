void setup()
{
    // 初始化串口通讯
    Serial.begin(9600);
}

void loop()
{
    // 发送模拟值
    Serial.println(analogRead(A0));
    // 等待模数转换器稳定下来
    delay(2);
}