#define ultra_trig 2
#define ultra_echo 3

void setup() {
    pinMode(ultra_trig, OUTPUT);
    pinMode(ultra_echo, INPUT);
    Serial.begin(9600);
}

void loop() {
    int duration;
    int distance;
    digitalWrite(ultra_trig, HIGH);
    delayMicroseconds(1000);
    digitalWrite(ultra_trig, LOW);
    duration = pulseIn(ultra_echo, HIGH);
    distance = (duration/2) /29.1;
    Serial.print(distance);
    Serial.println("cm");

    delay(500);
}