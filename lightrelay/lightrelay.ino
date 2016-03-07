// Function
//      Control LEDs and 2-way relay with Light sensor
// Library
//      https://github.com/PaulStoffregen/OneWire
//      https://github.com/milesburton/Arduino-Temperature-Control-Library
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
//      LED: Gnd->330Ohm->pinLedL->D2;
//           Gnd->330Ohm->pinLedM->D3;
//           Gnd->330Ohm->pinLedH->D4.
//      Relay module->Arduino: Gnd->Gnd(Arduino); In1->D5; In2->D6; Vcc->5V.
//      Relay module->Power switcher: NO1->AC PowerA1; NO2->AC PowerA2.
//      Relay module->Power switcher: COM1->COM2->AC PowerB.
// Tested
//      arduino uno/nano328
#include <OneWire.h>
#include <DallasTemperature.h>

#define BAUDRATE 9600

#define pinLedL 2
#define pinLedM 3
#define pinLedH 4
#define pinRelay1 5  // HIGH ON
#define pinRelay2 6  // HIGH OFF
#define pinDS18B20 7 // one wire bus
#define pinReadLight A0

int Relay1 = LOW;
int Relay2 = LOW;
int lightValue = 0;
float temperature = 0.0;

OneWire oneWire(pinDS18B20);
DallasTemperature sensors(&oneWire);

void ctrlLed(int l, int m, int h) {
    digitalWrite(pinLedL, l);
    digitalWrite(pinLedM, m);
    digitalWrite(pinLedH, h);
}

void setup()
{
    Serial.begin(BAUDRATE);

    pinMode(pinLedL, OUTPUT);
    pinMode(pinLedM, OUTPUT);
    pinMode(pinLedH, OUTPUT);
    pinMode(pinRelay1, OUTPUT);
    pinMode(pinRelay2, OUTPUT);

    ctrlLed(LOW, LOW, LOW);
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

    Serial.print(temperature, lightValue);

    if (lightValue > 950) {
        ctrlLed(HIGH, HIGH, HIGH);

        if (Relay2 == HIGH) {
            Relay2 = LOW;
            digitalWrite(pinRelay2, Relay2);
//            Serial.println("[X] Relay 2 Off");
        }
    } else if (lightValue > 800) {
        ctrlLed(HIGH, HIGH, LOW);
    } else if (lightValue > 700) {
        ctrlLed(HIGH, LOW, LOW);
    } else {
        ctrlLed(LOW, LOW, LOW);

        if (Relay2 == LOW) {
            Relay2 = HIGH;
            digitalWrite(pinRelay2, Relay2);
//            Serial.println("[X] Relay 2 On");
        }
    }

    delay(1000);
}
