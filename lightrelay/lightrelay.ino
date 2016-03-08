// Function
//      Control LEDs and 2-way relay with Light sensor
// Library
//      cd /usr/share/arduino/libraries/
//      git clone https://github.com/PaulStoffregen/OneWire
//      git clone https://github.com/milesburton/Arduino-Temperature-Control-Library DallasTemperature
//
// Reference
//      http://www.geek-workshop.com/thread-1340-1-1.html
//      http://www.instructables.com/id/DS18B20-temperature-probe-with-LCD/?ALLSTEPS
// Connection
//      Light Sensor: Gnd->10KOhm->LS->5V;
//                               |->A0.
//      Temperature Sensor - DS18B20: Vcc->5V(Arduino);
//                                    Gnd->Gnd(Arduino);
//                                    Data->4.7KOhm->5V(Arduino);
//                                        |->D7.
//      LED: Gnd->330Ohm->pinLedB->D2;
//           Gnd->330Ohm->pinLedG->D3;
//           Gnd->330Ohm->pinLedR->D4.
//      Relay module->Arduino: Gnd->Gnd(Arduino); In1->D5; In2->D6; Vcc->5V.
//      Relay module->Power switcher: NO1->AC PowerA1; NO2->AC PowerA2.
//      Relay module->Power switcher: COM1->COM2->AC PowerB.
// Tested
//      arduino uno/nano328
#include <OneWire.h>
#include <DallasTemperature.h>

#define BAUDRATE 9600

#define pinLedB 2
#define pinLedG 3
#define pinLedR 4
#define pinRelay1 5
#define pinRelay2 6
#define pinDS18B20 7 // one wire bus
#define pinReadLight A0

int Relay1 = LOW; // ON
int Relay2 = LOW; // ON
int lightValue = 0;
float temperature = 0.0;

OneWire oneWire(pinDS18B20);
DallasTemperature sensors(&oneWire);

void setup()
{
    Serial.begin(BAUDRATE);
    Serial.println("[INFO] Init");

    pinMode(pinLedR, OUTPUT);
    pinMode(pinLedG, OUTPUT);
    pinMode(pinLedB, OUTPUT);
    pinMode(pinRelay1, OUTPUT);
    pinMode(pinRelay2, OUTPUT);

    digitalWrite(pinLedR, LOW);
    digitalWrite(pinLedG, LOW);
    digitalWrite(pinLedB, LOW);

    digitalWrite(pinRelay1, Relay1);
    digitalWrite(pinRelay2, Relay2);
    Serial.println("[INFO] Relay 1 On");
    Serial.println("[INFO] Relay 2 On");

    // Start up the library
    sensors.begin();

    delay(500);
}

void loop()
{
    // Get light sensor data
    lightValue = analogRead(pinReadLight);

    // Control relay-1 with light sensor
    // On: >950, Off: <700
    if (lightValue > 950 && Relay1 == LOW) {
        Relay1 = HIGH;
        digitalWrite(pinRelay1, Relay1);
        digitalWrite(pinLedG, LOW);
        Serial.println("[INFO] Light above 950, cut off oxygen supply");
    } else if (lightValue < 700 && Relay1 == HIGH) {
        Relay1 = LOW;
        digitalWrite(pinRelay1, Relay1);
        digitalWrite(pinLedG, HIGH);
        Serial.println("[INFO] Light below 700, turn on oxygen supply");
    }

    // Send get-temperatures command
    sensors.requestTemperatures();
    // Get temperature sensor data at device 0
    temperature = sensors.getTempCByIndex(0);

    // Control relay-2 with temperature
    // On: <30.0, Off: >35.0
    if (temperature > 35.0 && Relay2 == LOW) {
        Relay2 = HIGH;
        digitalWrite(pinRelay2, Relay2);
        digitalWrite(pinLedR, HIGH);
        Serial.println("[INFO] Temperature above 35 degree, cut off power");
    } else if (temperature < 20.0 && Relay2 == HIGH) {
        Relay2 = LOW;
        digitalWrite(pinRelay2, Relay2);
        digitalWrite(pinLedR, LOW);
        Serial.println("[INFO] Temperature below 20 degree, turn on power");
    }

    // Write sensor data to serial port
    Serial.print(temperature);
    Serial.write(' ');
    Serial.println(lightValue);

    delay(1000);
}
