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
#define pinRelay1 5  // On when daylight >950, Off <700
#define pinRelay2 6  // On when temperature <30.0, Off >35.0
#define pinDS18B20 7 // one wire bus
#define pinReadLight A0

int Relay1 = LOW;
int Relay2 = LOW;
int lightValue = 0;
float temperature = 0.0;

OneWire oneWire(pinDS18B20);
DallasTemperature sensors(&oneWire);

void setup()
{
    Serial.begin(BAUDRATE);

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

    // Start up the library
    sensors.begin();
}

void loop()
{
    // Send the command to get temperatures
    sensors.requestTemperatures();
    // Why "byIndex"?
    // You can have more than one IC on the same bus.
    // 0 refers to the first IC on the wire
    temperature = sensors.getTempCByIndex(0);
    lightValue = analogRead(pinReadLight);

    Serial.print(temperature);
    Serial.write(' ');
    Serial.println(lightValue);

    if (lightValue > 950 && Relay1 == HIGH) {
        // Relay 1 Off
        Relay1 = LOW;
        digitalWrite(pinRelay1, Relay1);
        digitalWrite(pinLedG, HIGH);
    } else if (lightValue < 700 && Relay1 == LOW) {
        // Relay 1 On
        Relay1 = HIGH;
        digitalWrite(pinRelay1, Relay1);
        digitalWrite(pinLedG, LOW);
    }

    if (temperature > 35.0 && Relay2 == HIGH) {
        Relay2 = LOW;
        digitalWrite(pinRelay2, Relay2);
        digitalWrite(pinLedR, HIGH);
    } else if (temperature < 30.0 && Relay2 == LOW) {
        Relay2 = HIGH;
        digitalWrite(pinRelay2, Relay2);
        digitalWrite(pinLedR, LOW);
    }

    delay(1000);
}
