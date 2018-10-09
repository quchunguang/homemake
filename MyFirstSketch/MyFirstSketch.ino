void setup()
{
    Serial.begin(9600);
    pinMode(LED_BUILTIN, OUTPUT);
}

void loop()
{
    digitalWrite(LED_BUILTIN, HIGH);
    Serial.print('1');
    delay(200);
    digitalWrite(LED_BUILTIN, LOW);
    Serial.print('0');
    delay(200);
}
